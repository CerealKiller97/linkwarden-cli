package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGetV2RendersStatsAndLinks(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v2/dashboard", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"message":"Success","data":{"links":[{"id":1,"name":"Go","url":"https://go.dev"}],"numberOfPinnedLinks":2,"numberOfTags":5}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := GetV2(nil)
	assert.NoError(cmd.Execute())
}

func TestGetV2PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := GetV2(nil)
	assert.Error(cmd.Execute())
}
