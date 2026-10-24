package collections

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestCreate_UsesParentIDFlagWithoutListing(t *testing.T) {
	assert := require.New(t)

	var listCalled bool
	var gotBody linkwarden.CreateCollectionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			listCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/collections":
			assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":10,"name":"New Collection"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Create()
	cmd.SetArgs([]string{"--name", "New Collection", "--parent-id", "5"})
	assert.NoError(cmd.Execute())
	assert.False(listCalled)
	assert.Equal("New Collection", gotBody.Name)
	assert.Equal(5, gotBody.ParentID)
}

func TestCreate_LooksUpUnorganizedWhenParentIDOmitted(t *testing.T) {
	assert := require.New(t)

	var gotBody linkwarden.CreateCollectionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/collections":
			assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":10,"name":"New Collection"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())

	cmd := Create()
	cmd.SetArgs([]string{"--name", "New Collection"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, gotBody.ParentID)
}

func TestCreate_DuplicateNameErrorsWithoutCreating(t *testing.T) {
	assert := require.New(t)

	var createCalled bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			w.WriteHeader(http.StatusOK)
			// The duplicate comes after Unorganized, which the old loop stopped at.
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"},{"id":4,"name":"Reading"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/collections":
			createCalled = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":10,"name":"Reading"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	t.Setenv("XDG_STATE_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())

	cmd := Create()
	cmd.SetArgs([]string{"--name", "Reading", "--parent-id", "0"})
	err := cmd.Execute()
	assert.ErrorIs(err, errCollectionExists)
	assert.False(createCalled)
}
