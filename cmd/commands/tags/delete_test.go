package tags

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestDeleteUsesIDFlag(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/tags/3", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Tag deleted."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete(nil)
	cmd.SetArgs([]string{"--id", "3", "--yes"})
	assert.NoError(cmd.Execute())
}

func TestDeletePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"forbidden"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete(nil)
	cmd.SetArgs([]string{"--id", "3", "--yes"})
	assert.Error(cmd.Execute())
}

func TestDeleteRequiresYesWithoutTerminal(t *testing.T) {
	assert := require.New(t)

	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete(nil)
	cmd.SetArgs([]string{"--id", "3"})
	err := cmd.Execute()
	assert.ErrorIs(err, utils.ErrConfirmationRequired)
	assert.Zero(requests, "nothing may reach the API without --yes")
}
