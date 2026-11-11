package links

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestListRendersLinks(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/search", r.URL.Path)
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Success","data":{"links":[{"id":1,"name":"Example","url":"https://example.com"}]}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := List(nil)
	cmd.SetArgs([]string{"--collection-id", "2", "--pinned-only"})
	assert.NoError(cmd.Execute())
	assert.Equal("collectionId=2&pinnedOnly=true", gotQuery)
}

func TestListPropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := List(nil)
	cmd.SetArgs([]string{})
	assert.Error(cmd.Execute())
}
