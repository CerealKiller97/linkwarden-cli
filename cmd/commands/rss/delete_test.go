package rss

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
		assert.Equal("/api/v1/rss/5", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"RSS subscription deleted."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete(nil)
	cmd.SetArgs([]string{"--id", "5", "--yes"})
	assert.NoError(cmd.Execute())
}

func TestDeleteInvalidIDErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Delete(nil)
	cmd.SetArgs([]string{"--id", "not-a-number", "--yes"})
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
	cmd.SetArgs([]string{"--id", "5"})
	err := cmd.Execute()
	assert.ErrorIs(err, utils.ErrConfirmationRequired)
	assert.Zero(requests, "nothing may reach the API without --yes")
}
