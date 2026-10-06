package collections

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestDelete_UsesIDFlag(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/collections/4", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Collection deleted."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", "4"})
	assert.NoError(cmd.Execute())
}

func TestDelete_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"Cannot delete."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", "4"})
	err := cmd.Execute()
	assert.Error(err)
}
