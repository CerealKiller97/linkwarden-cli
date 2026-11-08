package tokens

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestRevokeFetchesTokensBeforePrompting(t *testing.T) {
	assert := require.New(t)

	var listCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"CI token"}]}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Revoke()
	err := cmd.Execute()
	assert.True(listCalled)
	// The picker needs an interactive terminal; in a non-interactive test run it
	// errors out immediately, which is expected here — this test only asserts that
	// the token list is fetched before the interactive picker runs.
	assert.Error(err)
}

func TestRevokePropagatesListError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Revoke()
	assert.Error(cmd.Execute())
}

func TestRevokeHasYesFlag(t *testing.T) {
	assert := require.New(t)

	flag := Revoke().Flags().Lookup("yes")
	assert.NotNil(flag)
	assert.Equal("y", flag.Shorthand)
}
