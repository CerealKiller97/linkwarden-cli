package links

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreate_UploadsLinksFromFile(t *testing.T) {
	assert := require.New(t)

	var createCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/links":
			createCount++
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":1,"name":"x","url":"https://example.com"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	file := filepath.Join(t.TempDir(), "links.json")
	assert.NoError(os.WriteFile(file, []byte(`[{"url":"https://example.com","title":"Example"}]`), 0o644))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, createCount)
}

func TestCreate_EmptyFileErrors(t *testing.T) {
	assert := require.New(t)

	file := filepath.Join(t.TempDir(), "links.json")
	assert.NoError(os.WriteFile(file, []byte(`[]`), 0o644))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})
	err := cmd.Execute()
	assert.ErrorIs(err, errFileEmpty)
}
