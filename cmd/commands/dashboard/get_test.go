package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGet_RendersLinks(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/dashboard", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Go","url":"https://go.dev"}]}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Get()
	assert.NoError(cmd.Execute())
}

func TestGet_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Get()
	assert.Error(cmd.Execute())
}
