package tokens

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils/prompttest"
	"github.com/stretchr/testify/require"
)

// fakeTokensAPI serves the token endpoints and records what the CLI sent.
type fakeTokensAPI struct {
	mu        sync.Mutex
	created   []map[string]any
	revoked   []string
	createErr bool
	list      string
}

func (f *fakeTokensAPI) start(t *testing.T) {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/tokens", func(w http.ResponseWriter, r *http.Request) {
		if f.createErr {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"response":"boom"}`))
			return
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.mu.Lock()
		f.created = append(f.created, body)
		f.mu.Unlock()
		_, _ = w.Write([]byte(`{"response":{"secretKey":"new-secret-token"}}`))
	})
	mux.HandleFunc("GET /api/v1/tokens", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(f.list))
	})
	mux.HandleFunc("DELETE /api/v1/tokens/{id}", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.revoked = append(f.revoked, r.PathValue("id"))
		f.mu.Unlock()
		_, _ = w.Write([]byte(`{"response":"ok"}`))
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	utils.SetVars("test-token", srv.URL)
}

// homeWithShell points HOME (USERPROFILE on Windows) and SHELL at a temporary home so nothing touches the
// real shell config, and returns the zsh config path the CLI would write.
func homeWithShell(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // os.UserHomeDir reads this on Windows
	t.Setenv("SHELL", "/bin/zsh")

	return filepath.Join(home, ".zshenv")
}

// captureStdout returns what fn printed to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()

	fn()
	require.NoError(t, w.Close())

	return <-done
}

func skipUnlessShellConfigOffered(t *testing.T) {
	t.Helper()

	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("tokens generate only offers to edit shell config on Linux and macOS")
	}
}

func TestGenerateCreatesTokenAndWritesShellConfig(t *testing.T) {
	skipUnlessShellConfigOffered(t)
	assert := require.New(t)

	api := &fakeTokensAPI{}
	api.start(t)
	zshenv := homeWithShell(t)

	prompttest.Use(t,
		prompttest.Answer{Label: "Token name", Value: "laptop"},
		prompttest.Answer{Label: "Token expiration", Value: "30 days"},
		prompttest.Answer{Label: "set the token in your shell config", Value: "Yes"},
	)

	assert.NoError(Generate(nil).Execute())

	assert.Equal([]map[string]any{{"name": "laptop", "expires": float64(ThirtyDays)}}, api.created)

	data, err := os.ReadFile(zshenv)
	assert.NoError(err)
	assert.Contains(string(data), `export LINKWARDEN_TOKEN="new-secret-token"`)

	info, err := os.Stat(zshenv)
	assert.NoError(err)
	assert.Equal(os.FileMode(0o600), info.Mode().Perm(), "the shell config now holds a secret")
}

func TestGenerateDeclinedPrintsManualStepsAndLeavesShellConfigAlone(t *testing.T) {
	skipUnlessShellConfigOffered(t)
	assert := require.New(t)

	api := &fakeTokensAPI{}
	api.start(t)
	zshenv := homeWithShell(t)
	t.Setenv("NO_COLOR", "1")

	prompttest.Use(t,
		prompttest.Answer{Label: "Token name", Value: "laptop"},
		prompttest.Answer{Label: "Token expiration", Value: "no expiration"},
		prompttest.Answer{Label: "set the token in your shell config", Value: "No"},
	)

	var err error
	out := captureStdout(t, func() { err = Generate(nil).Execute() })
	assert.NoError(err)

	assert.Equal(float64(NoExpiry), api.created[0]["expires"])
	assert.Contains(out, `export LINKWARDEN_TOKEN="new-secret-token"`)
	assert.Contains(out, "~/.zshenv")
	assert.NoFileExists(zshenv, "declining must not touch the shell config")
}

func TestGenerateAPIErrorSkipsShellConfig(t *testing.T) {
	assert := require.New(t)

	api := &fakeTokensAPI{createErr: true}
	api.start(t)
	zshenv := homeWithShell(t)

	// No answer for the shell-config question: asking it would fail the test.
	prompttest.Use(t,
		prompttest.Answer{Label: "Token name", Value: "laptop"},
		prompttest.Answer{Label: "Token expiration", Value: "7 days"},
	)

	assert.Error(Generate(nil).Execute())
	assert.NoFileExists(zshenv)
}

func TestGenerateRejectsBlankTokenName(t *testing.T) {
	assert := require.New(t)

	api := &fakeTokensAPI{}
	api.start(t)
	homeWithShell(t)

	var validationErr error
	restore := utils.SetPrompter(validatingPrompter{onInput: func(validate func(string) error) {
		validationErr = validate("   ")
	}})
	defer restore()

	_ = Generate(nil).Execute()
	assert.EqualError(validationErr, "token name is required")
	assert.Empty(api.created)
}

// validatingPrompter hands the name prompt's validation to onInput, then aborts.
type validatingPrompter struct {
	utils.Prompter
	onInput func(validate func(string) error)
}

func (p validatingPrompter) Input(_ string, validate func(string) error) (string, error) {
	p.onInput(validate)
	return "", io.EOF
}

func TestRevokeRevokesTheChosenTokenWhenNamesRepeat(t *testing.T) {
	assert := require.New(t)

	api := &fakeTokensAPI{list: `{"response":[{"id":11,"name":"CI"},{"id":12,"name":"laptop"},{"id":13,"name":"CI"}]}`}
	api.start(t)

	// Two tokens are named "CI": picking the second must revoke id 13, not the
	// first "CI" a lookup by name would find.
	prompttest.Use(t,
		prompttest.Pick("Select token to revoke", 2),
		prompttest.Answer{Label: `Revoke token "CI"?`, Value: "Yes"},
	)

	assert.NoError(Revoke(nil).Execute())
	assert.Equal([]string{"13"}, api.revoked)
}

func TestRevokeAnsweringNoRevokesNothing(t *testing.T) {
	assert := require.New(t)

	api := &fakeTokensAPI{list: `{"response":[{"id":12,"name":"laptop"}]}`}
	api.start(t)

	prompttest.Use(t,
		prompttest.Answer{Label: "Select token to revoke", Value: "laptop"},
		prompttest.Answer{Label: "Revoke token", Value: "No"},
	)

	assert.NoError(Revoke(nil).Execute())
	assert.Empty(api.revoked)
}

func TestRevokeYesFlagSkipsConfirmation(t *testing.T) {
	assert := require.New(t)

	api := &fakeTokensAPI{list: `{"response":[{"id":12,"name":"laptop"}]}`}
	api.start(t)

	// Only the picker is answered: a confirmation prompt would fail the test.
	prompttest.Use(t, prompttest.Answer{Label: "Select token to revoke", Value: "laptop"})

	cmd := Revoke(nil)
	cmd.SetArgs([]string{"--yes"})
	assert.NoError(cmd.Execute())
	assert.Equal([]string{"12"}, api.revoked)
}
