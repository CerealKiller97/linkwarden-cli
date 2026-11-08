package tags

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestListRendersTags(t *testing.T) {
	assert := require.New(t)

	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/tags", r.URL.Path)
		gotQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Success","success":true,"data":{"tags":[{"id":1,"name":"Technology"}]}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := List()
	cmd.SetArgs([]string{"--search", "tech"})
	assert.NoError(cmd.Execute())
	assert.Equal("search=tech", gotQuery)
}

func TestListPropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := List()
	assert.Error(cmd.Execute())
}
