package links

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestArchiveUsesIDFlag(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/links/8/archive", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Link archived."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Archive(nil)
	cmd.SetArgs([]string{"--id", "8"})
	assert.NoError(cmd.Execute())
}

func TestArchivePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Archive(nil)
	cmd.SetArgs([]string{"--id", "8"})
	assert.Error(cmd.Execute())
}
