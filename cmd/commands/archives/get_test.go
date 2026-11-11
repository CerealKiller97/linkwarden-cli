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

func TestGetSavesToOutputFile(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/archives/5", r.URL.Path)
		assert.Equal("format=2", r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("%PDF-fake"))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	outputPath := filepath.Join(t.TempDir(), "archive.pdf")

	cmd := Get(nil)
	cmd.SetArgs([]string{"--link-id", "5", "--format", "pdf", "--output", outputPath})
	assert.NoError(cmd.Execute())

	data, err := os.ReadFile(outputPath)
	assert.NoError(err)
	assert.Equal("%PDF-fake", string(data))
}

func TestGetWithoutOutputPrintsSummary(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("fake-png"))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Get(nil)
	cmd.SetArgs([]string{"--link-id", "5"})
	assert.NoError(cmd.Execute())
}

func TestGetPropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"response":"Link not found."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Get(nil)
	cmd.SetArgs([]string{"--link-id", "999"})
	assert.Error(cmd.Execute())
}
