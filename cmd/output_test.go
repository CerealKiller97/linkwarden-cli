package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

// runRoot runs the CLI with args against server and returns what it printed
// on stdout.
func runRoot(t *testing.T, server *httptest.Server, args ...string) (string, error) {
	t.Helper()

	t.Setenv("LINKWARDEN_TOKEN", "test-token")
	t.Setenv("LINKWARDEN_URL", server.URL)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Cleanup(func() { utils.SetOutputMode(utils.OutputTable) })

	r, w, err := os.Pipe()
	require.NoError(t, err)

	stdout := os.Stdout
	os.Stdout = w
	root := Root()
	root.SetArgs(args)
	runErr := root.Execute()
	os.Stdout = stdout
	require.NoError(t, w.Close())

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(out), runErr
}

func tagsServer(t *testing.T) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/v1/tags" {
			_, _ = w.Write([]byte(`{"message":"Success","success":true,"data":{"tags":[{"id":3,"name":"go"},{"id":5,"name":"k8s"}]}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	return server
}

func TestJSONFlagPrintsJSON(t *testing.T) {
	assert := require.New(t)

	out, err := runRoot(t, tagsServer(t), "tags", "list", "--json")
	assert.NoError(err)

	var tags []map[string]any
	assert.NoError(json.Unmarshal([]byte(out), &tags), "stdout is valid JSON: %q", out)
	assert.Len(tags, 2)
	assert.Equal("go", tags[0]["name"])
	assert.InDelta(3, tags[0]["id"], 0)
}

func TestQuietFlagPrintsIDs(t *testing.T) {
	assert := require.New(t)

	out, err := runRoot(t, tagsServer(t), "tags", "list", "-q")
	assert.NoError(err)
	assert.Equal("3\n5\n", out)
}

func TestJSONAndQuietTogetherError(t *testing.T) {
	_, err := runRoot(t, tagsServer(t), "tags", "list", "--json", "--quiet")
	require.ErrorIs(t, err, errJSONAndQuiet)
}
