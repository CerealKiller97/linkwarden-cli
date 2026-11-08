package links

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestUpdateSendsFlagsAsRequestBody(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/links/9", r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{"response":{"id":9,"name":"Old","url":"https://example.com","collection":{"id":1,"ownerId":1}}}`))
		case http.MethodPut:
			assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
			_, _ = w.Write([]byte(`{"response":{"id":9,"name":"Renamed","url":"https://example.com"}}`))
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "9", "--name", "Renamed"})
	assert.NoError(cmd.Execute())
	assert.Equal("Renamed", gotBody["name"])
	assert.Equal("https://example.com", gotBody["url"])
}

func TestUpdatePropagatesAPIError(t *testing.T) {
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
