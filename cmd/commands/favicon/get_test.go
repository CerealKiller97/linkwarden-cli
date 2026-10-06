package favicon

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGet_SavesToOutputFile(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/getFavicon", r.URL.Path)
		assert.Equal("https://example.com", r.URL.Query().Get("url"))
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fake-icon"))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	outputPath := filepath.Join(t.TempDir(), "favicon.ico")

	cmd := Get()
	cmd.SetArgs([]string{"--url", "https://example.com", "--output", outputPath})
	assert.NoError(cmd.Execute())

	data, err := os.ReadFile(outputPath)
	assert.NoError(err)
	assert.Equal("fake-icon", string(data))
}

func TestGet_WithoutOutputPrintsSummary(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fake-png"))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Get()
	cmd.SetArgs([]string{"--url", "https://example.com"})
	assert.NoError(cmd.Execute())
}

func TestGet_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid url"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Get()
	cmd.SetArgs([]string{"--url", "not-a-url"})
	assert.Error(cmd.Execute())
}
