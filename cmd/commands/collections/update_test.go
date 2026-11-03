package collections

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

// updateServer serves the GET the SDK issues before the PUT and records the PUT body.
func updateServer(t *testing.T, current string, gotBody *map[string]any) *httptest.Server {
	assert := require.New(t)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/collections/3", r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			_, _ = w.Write([]byte(`{"response":` + current + `}`))
		case http.MethodPut:
			assert.NoError(json.NewDecoder(r.Body).Decode(gotBody))
			_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Renamed"}}`))
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	}))
}

func TestUpdate_SendsFlagsAsRequestBody(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	server := updateServer(t, `{"id":3,"name":"Old","isPublic":false}`, &gotBody)
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "3", "--name", "Renamed", "--public"})
	assert.NoError(cmd.Execute())
	assert.Equal("Renamed", gotBody["name"])
	assert.Equal(true, gotBody["isPublic"])
}

func TestUpdate_PublicFalseMakesPrivate(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	server := updateServer(t, `{"id":3,"name":"Old","isPublic":true}`, &gotBody)
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "3", "--public=false"})
	assert.NoError(cmd.Execute())
	assert.Equal("Old", gotBody["name"])
	assert.Equal(false, gotBody["isPublic"])
}

func TestUpdate_OmittedPublicKeepsCurrent(t *testing.T) {
	assert := require.New(t)

	var gotBody map[string]any
	server := updateServer(t, `{"id":3,"name":"Old","isPublic":true}`, &gotBody)
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "3", "--name", "Renamed"})
	assert.NoError(cmd.Execute())
	assert.Equal(true, gotBody["isPublic"])
}

func TestUpdate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Collection not found."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "3"})
	assert.Error(cmd.Execute())
}
