package migrations

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestExportWritesOutputFile(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/migration", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"username":"bob","collections":[]}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	outputPath := filepath.Join(t.TempDir(), "export.json")

	cmd := Export(nil)
	cmd.SetArgs([]string{"--output", outputPath})
	assert.NoError(cmd.Execute())

	data, err := os.ReadFile(outputPath)
	assert.NoError(err)
	assert.JSONEq(`{"username":"bob","collections":[]}`, string(data))
}

func TestExportPropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Export(nil)
	cmd.SetArgs([]string{"--output", filepath.Join(t.TempDir(), "export.json")})
	assert.Error(cmd.Execute())
}
