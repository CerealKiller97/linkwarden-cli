package search

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestExecuteUsesQueryFlags(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/search", r.URL.Path)
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Success","data":{"links":[{"id":1,"name":"Go","url":"https://go.dev"}]}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Command()
	cmd.SetArgs([]string{"--query", "go", "--sort", "name-asc", "--cursor", "3"})
	assert.NoError(cmd.Execute())
	assert.Equal("cursor=3&searchQueryString=go&sort=2", gotQuery)
}

func TestExecuteOmitsUnsetSort(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Success","data":{"links":[]}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Command()
	cmd.SetArgs([]string{})
	assert.NoError(cmd.Execute())
	assert.Equal("", gotQuery)
}

func TestExecutePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Command()
	assert.Error(cmd.Execute())
}
