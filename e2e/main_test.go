//go:build e2e

package e2e

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/stretchr/testify/require"
)

var (
	// binary is the CLI built from this checkout.
	binary string
	// stateDir isolates the CLI's cached state (e.g. the Unorganized collection ID)
	// from the developer's own.
	stateDir string
	// api is used to set up fixtures, verify what the CLI did, and clean up.
	api *linkwarden.Client
)

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	instanceURL, token := os.Getenv("LINKWARDEN_URL"), os.Getenv("LINKWARDEN_TOKEN")
	if instanceURL == "" || token == "" {
		fmt.Fprintln(os.Stderr, "e2e: LINKWARDEN_URL and LINKWARDEN_TOKEN must be set; see scripts/dev-linkwarden.sh")
		return 1
	}

	u, err := url.Parse(instanceURL)
	if err != nil || (u.Hostname() != "localhost" && u.Hostname() != "127.0.0.1") {
		fmt.Fprintf(os.Stderr, "e2e: refusing to run against %q; only a local instance is allowed\n", instanceURL)
		return 1
	}

	dir, err := os.MkdirTemp("", "linkwarden-e2e")
	if err != nil {
		fmt.Fprintln(os.Stderr, "e2e:", err)
		return 1
	}
	defer func() { _ = os.RemoveAll(dir) }()

	binary = filepath.Join(dir, "linkwarden")
	build := exec.Command("go", "build", "-o", binary, "..")
	build.Stdout, build.Stderr = os.Stderr, os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "e2e: build CLI:", err)
		return 1
	}

	stateDir = filepath.Join(dir, "state")
	api = linkwarden.New(instanceURL, token)

	return m.Run()
}

// cli runs the CLI non-interactively with extra environment variables and
// returns its combined output and exit code.
func cliWithEnv(t *testing.T, env []string, args ...string) (string, int) {
	t.Helper()

	cmd := exec.CommandContext(t.Context(), binary, append([]string{"--non-interactive"}, args...)...)
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "XDG_DATA_HOME="+stateDir, "XDG_STATE_HOME="+stateDir)
	cmd.Env = append(cmd.Env, env...)

	out, err := cmd.CombinedOutput()
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		return string(out), exitErr.ExitCode()
	}
	require.NoError(t, err, "run linkwarden %s", strings.Join(args, " "))

	return string(out), 0
}

// mustCLI runs the CLI and fails the test unless it exits 0.
func mustCLI(t *testing.T, args ...string) string {
	t.Helper()

	out, code := cliWithEnv(t, nil, args...)
	require.Zero(t, code, "linkwarden %s exited %d:\n%s", strings.Join(args, " "), code, out)

	return out
}

// failCLI runs the CLI and fails the test unless it exits non-zero.
func failCLI(t *testing.T, args ...string) string {
	t.Helper()

	out, code := cliWithEnv(t, nil, args...)
	require.NotZero(t, code, "linkwarden %s unexpectedly succeeded:\n%s", strings.Join(args, " "), out)

	return out
}

var loggedIDPattern = regexp.MustCompile(`\bid=(\d+)`)

// loggedID extracts the ID from a success log line such as
// "Collection created successfully id=12".
func loggedID(t *testing.T, out string) int {
	t.Helper()

	match := loggedIDPattern.FindStringSubmatch(out)
	require.NotNil(t, match, "no id=N in output:\n%s", out)

	id, err := strconv.Atoi(match[1])
	require.NoError(t, err)

	return id
}

// unique returns a short name that won't collide with other runs or tests.
func unique(base string) string {
	return fmt.Sprintf("e2e-%s-%s", base, strconv.FormatInt(time.Now().UnixNano(), 36))
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

// cleanup registers fn to run after the test with a context that outlives
// t.Context(), which is already canceled by the time cleanups run.
func cleanup(t *testing.T, fn func(ctx context.Context) error) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = fn(ctx)
	})
}

// newCollection creates a top-level collection that is deleted after the test,
// together with any links still in it.
func newCollection(t *testing.T) linkwarden.Collection {
	t.Helper()

	collection, err := api.Collections.Create(t.Context(), linkwarden.CreateCollectionRequest{Name: unique("collection")})
	require.NoError(t, err)
	cleanup(t, func(ctx context.Context) error { return api.Collections.Delete(ctx, collection.ID) })

	return *collection
}

// newLink creates a link in collection. It goes away with the collection.
func newLink(t *testing.T, collection linkwarden.Collection, req linkwarden.CreateLinkRequest) linkwarden.Link {
	t.Helper()

	if req.Name == "" {
		req.Name = unique("link")
	}
	if req.URL == "" {
		req.URL = "https://example.com/" + req.Name
	}
	req.Type = "url"
	req.Collection = linkwarden.Collection{ID: collection.ID}

	link, err := api.Links.Create(t.Context(), req)
	require.NoError(t, err)

	return *link
}

// findTag looks a tag up by exact name, failing the test if it doesn't exist.
func findTag(t *testing.T, name string) linkwarden.TagInfo {
	t.Helper()

	tags, err := api.Tags.List(t.Context(), linkwarden.ListTagsOptions{Search: name})
	require.NoError(t, err)

	for _, tag := range tags {
		if tag.Name == name {
			return tag
		}
	}
	require.FailNow(t, "tag not found", "no tag named %q", name)

	return linkwarden.TagInfo{}
}

// eventually retries check until it returns true or the timeout passes, for
// work the server finishes asynchronously (search indexing, imports).
func eventually(t *testing.T, timeout time.Duration, msg string, check func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for !check() {
		if time.Now().After(deadline) {
			require.FailNow(t, "timed out", msg)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

// writeFile writes contents to a file in a per-test temp dir and returns its path.
func writeFile(t *testing.T, name, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	return path
}
