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

func TestUpdate_MergesChangedFlagsIntoCurrentCollection(t *testing.T) {
	assert := require.New(t)

	var gotBody linkwarden.UpdateCollectionRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/collections/3", r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Old","description":"keep me","color":"#123456","parentId":7,
				"members":[{"userId":5,"canCreate":true,"canUpdate":false,"canDelete":false}]}}`))
		case http.MethodPut:
			assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Renamed"}}`))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "3", "--name", "Renamed", "--public"})
	assert.NoError(cmd.Execute())
	assert.Equal(3, gotBody.ID)
	assert.Equal("Renamed", gotBody.Name)
	assert.True(*gotBody.IsPublic)
	assert.Equal("keep me", gotBody.Description)
	assert.Equal("#123456", gotBody.Color)
	assert.Equal(7, *gotBody.ParentID)
	assert.Len(gotBody.Members, 1)
	assert.Equal(5, gotBody.Members[0].UserID)
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
