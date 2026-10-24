//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/stretchr/testify/require"
)

func TestConfigGet(t *testing.T) {
	require.Regexp(t, `INSTANCE_VERSION\s+│\s+v\d+\.\d+`, mustCLI(t, "config", "get"))
}

func TestTokensList(t *testing.T) {
	tokens, err := api.Tokens.List(t.Context())
	require.NoError(t, err)
	require.NotEmpty(t, tokens)

	require.Contains(t, mustCLI(t, "tokens", "list"), tokens[0].Name)
}

func TestDashboard(t *testing.T) {
	mustCLI(t, "dashboard", "get")
	mustCLI(t, "dashboard", "get-v2")

	// order 0 used to be dropped from the request, which the API answers with a 500.
	layout := writeFile(t, "layout.json",
		`[{"type":"STATS","enabled":true,"order":0},{"type":"RECENT_LINKS","enabled":true,"order":1}]`)
	mustCLI(t, "dashboard", "update-v2", "--file", layout)
}

func TestArchives_UploadAndGet(t *testing.T) {
	assert := require.New(t)

	link := newLink(t, newCollection(t), linkwarden.CreateLinkRequest{})
	page := writeFile(t, "page.html", "<html><body><h1>"+link.Name+"</h1></body></html>")

	mustCLI(t, "archives", "upload", "--link-id", itoa(link.ID), "--format", "4", "--file", page)

	// The server may replace the archive with its own capture of the URL at any
	// time, so only check that an HTML archive comes back.
	saved := filepath.Join(t.TempDir(), "archive.html")
	mustCLI(t, "archives", "get", "--link-id", itoa(link.ID), "--format", "4", "--output", saved)
	contents, err := os.ReadFile(saved)
	assert.NoError(err)
	assert.Contains(string(contents), "<html")
}

var tableIDPattern = regexp.MustCompile(`│\s+(\d+)\s+│`)

func TestArchives_UploadStandalone(t *testing.T) {
	assert := require.New(t)

	name := unique("standalone")
	page := writeFile(t, "page.html", "<html><body><h1>"+name+"</h1></body></html>")

	out := mustCLI(t, "archives", "upload-standalone", "--url", "https://example.com/"+name, "--format", "4", "--file", page)
	match := tableIDPattern.FindStringSubmatch(out)
	assert.NotNil(match, "no link ID in output:\n%s", out)
	id, err := strconv.Atoi(match[1])
	assert.NoError(err)
	cleanup(t, func(ctx context.Context) error { return api.Links.Delete(ctx, id) })

	link, err := api.Links.Get(t.Context(), id)
	assert.NoError(err)
	assert.Equal("https://example.com/"+name, link.URL)
}

func TestMigrations_Export(t *testing.T) {
	assert := require.New(t)

	collection := newCollection(t)
	path := filepath.Join(t.TempDir(), "export.json")
	mustCLI(t, "migrations", "export", "--output", path)

	contents, err := os.ReadFile(path)
	assert.NoError(err)

	var export struct {
		Collections []linkwarden.Collection `json:"collections"`
	}
	assert.NoError(json.Unmarshal(contents, &export))

	var names []string
	for _, c := range export.Collections {
		names = append(names, c.Name)
	}
	assert.Contains(names, collection.Name)
}

func TestMigrations_ImportBookmarksHTML(t *testing.T) {
	folder := unique("import")
	bookmarks := writeFile(t, "bookmarks.html", `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<DL><p>
<DT><H3>`+folder+`</H3>
<DL><p>
<DT><A HREF="https://example.com/`+folder+`">Imported</A>
</DL><p>
</DL><p>
`)

	mustCLI(t, "migrations", "import", "--file", bookmarks, "--format", "1")

	var imported *linkwarden.Collection
	eventually(t, 30*time.Second, "imported collection "+folder+" never appeared", func() bool {
		collections, err := api.Collections.List(t.Context())
		if err != nil {
			return false
		}
		for _, c := range collections {
			if c.Name == folder {
				imported = &c
				return true
			}
		}
		return false
	})
	cleanup(t, func(ctx context.Context) error { return api.Collections.Delete(ctx, imported.ID) })

	links, err := api.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: imported.ID})
	require.NoError(t, err)
	require.Len(t, links, 1)
	require.Equal(t, "https://example.com/"+folder, links[0].URL)
}

func TestErrors(t *testing.T) {
	assert := require.New(t)

	out := failCLI(t, "links", "delete")
	assert.Contains(out, "pass --id explicitly")

	out = failCLI(t, "collections", "create")
	assert.Contains(out, `required flag(s) "name" not set`)

	out, code := cliWithEnv(t, []string{"LINKWARDEN_TOKEN=not-a-real-token"}, "collections", "list")
	assert.NotZero(code, "an invalid token should fail:\n%s", out)
}
