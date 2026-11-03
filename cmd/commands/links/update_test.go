package links

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestUpdate_SendsFlagsAsRequestBody(t *testing.T) {
	assert := require.New(t)

	var gotBody linkwarden.UpdateLinkRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/links/9", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":9,"name":"Renamed","url":"https://example.com"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "9", "--name", "Renamed"})
	assert.NoError(cmd.Execute())
	assert.Equal("Renamed", gotBody.Name)
}

func TestUpdate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Link not found."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "9"})
	assert.Error(cmd.Execute())
}
