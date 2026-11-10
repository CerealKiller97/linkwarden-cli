package migrations

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestImportSendsFileContentsAndFormat(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		Format int    `json:"format"`
		Data   string `json:"data"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/migration", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Import started."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	filePath := filepath.Join(t.TempDir(), "bookmarks.html")
	require.NoError(t, os.WriteFile(filePath, []byte("<html>bookmarks</html>"), 0o600))

	cmd := Import()
	cmd.SetArgs([]string{"--file", filePath, "--format", "html"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, gotBody.Format)
	assert.Equal("<html>bookmarks</html>", gotBody.Data)
}

func TestImportMissingFileErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Import()
	cmd.SetArgs([]string{"--file", "/nonexistent/file.json", "--format", "linkwarden"})
	assert.Error(cmd.Execute())
}
