package archives

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func writeTempFile(t *testing.T, name, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	return path
}

func TestUploadSendsFile(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/archives/2", r.URL.Path)
		assert.Equal("format=0", r.URL.RawQuery)

		assert.NoError(r.ParseMultipartForm(1 << 20))
		file, header, err := r.FormFile("file")
		assert.NoError(err)
		defer file.Close()
		assert.Equal("screenshot.png", header.Filename)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":2,"name":"Example","url":"https://example.com"}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	filePath := writeTempFile(t, "screenshot.png", "fake-png-bytes")

	cmd := Upload()
	cmd.SetArgs([]string{"--link-id", "2", "--file", filePath, "--format", "png"})
	assert.NoError(cmd.Execute())
}

func TestUploadMissingFileErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Upload()
	cmd.SetArgs([]string{"--link-id", "2", "--file", "/nonexistent/file.png", "--format", "png"})
	assert.Error(cmd.Execute())
}

func TestUploadRejectsReadableFormat(t *testing.T) {
	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Upload()
	cmd.SetArgs([]string{"--link-id", "2", "--file", "page.json", "--format", "readable"})
	require.ErrorContains(t, cmd.Execute(), "must be one of png|jpeg|pdf|html")
}
