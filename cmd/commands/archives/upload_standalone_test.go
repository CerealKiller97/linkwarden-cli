package archives

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestUploadStandaloneSendsFileAndURL(t *testing.T) {
	assert := require.New(t)

	var gotURLField string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/archives", r.URL.Path)
		assert.Equal("format=4&preview=true", r.URL.RawQuery)

		assert.NoError(r.ParseMultipartForm(1 << 20))
		gotURLField = r.FormValue("url")
		file, _, err := r.FormFile("file")
		assert.NoError(err)
		defer file.Close()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":10,"name":"Standalone","url":"https://example.com/page"}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	filePath := writeTempFile(t, "page.html", "<html></html>")

	cmd := UploadStandalone()
	cmd.SetArgs([]string{"--file", filePath, "--url", "https://example.com/page", "--format", "html", "--preview"})
	assert.NoError(cmd.Execute())
	assert.Equal("https://example.com/page", gotURLField)
}

func TestUploadStandalonePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid file"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	filePath := writeTempFile(t, "page.html", "<html></html>")

	cmd := UploadStandalone()
	cmd.SetArgs([]string{"--file", filePath, "--url", "https://example.com/page", "--format", "html"})
	assert.Error(cmd.Execute())
}
