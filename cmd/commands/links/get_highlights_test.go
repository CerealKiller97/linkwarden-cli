package links

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGetHighlightsUsesIDFlag(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/links/6/highlights", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"linkId":6,"text":"quote","color":"yellow"}]}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := GetHighlights(nil)
	cmd.SetArgs([]string{"--id", "6"})
	assert.NoError(cmd.Execute())
}

func TestGetHighlightsPropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := GetHighlights(nil)
	cmd.SetArgs([]string{"--id", "6"})
	assert.Error(cmd.Execute())
}
