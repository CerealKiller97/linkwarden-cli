package collections

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestDeleteUsesIDFlag(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/collections/4", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Collection deleted."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete(nil)
	cmd.SetArgs([]string{"--id", "4", "--yes"})
	assert.NoError(cmd.Execute())
}

func TestDeletePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"Cannot delete."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Delete(nil)
	cmd.SetArgs([]string{"--id", "4", "--yes"})
	err := cmd.Execute()
	assert.Error(err)
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
	cmd.SetArgs([]string{"--id", "4"})
	err := cmd.Execute()
	assert.ErrorIs(err, utils.ErrConfirmationRequired)
	assert.Zero(requests, "nothing may reach the API without --yes")
}

func TestDeleteQuestionNamesLinkCount(t *testing.T) {
	assert := require.New(t)

	assert.Equal(`Delete collection "Go" and its 12 links? This can't be undone.`,
		deleteQuestion(&linkwarden.Collection{Name: "Go", Count: &linkwarden.CollectionCount{Links: 12}}))
	assert.Equal(`Delete collection "Go" and its 1 link? This can't be undone.`,
		deleteQuestion(&linkwarden.Collection{Name: "Go", Count: &linkwarden.CollectionCount{Links: 1}}))
	assert.Equal(`Delete collection "Go" and all its links? This can't be undone.`,
		deleteQuestion(&linkwarden.Collection{Name: "Go"}))
}
