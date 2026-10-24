package links

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestUpdate_MergesChangedFlagsIntoCurrentLink(t *testing.T) {
	assert := require.New(t)

	var gotBody linkwarden.UpdateLinkRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/links/9", r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":9,"name":"Old","url":"https://example.com","description":"keep me",
				"collection":{"id":2,"ownerId":1,"name":"Reading"},"tags":[{"id":4,"name":"go"}]}}`))
		case http.MethodPut:
			assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":9,"name":"Renamed","url":"https://example.com"}}`))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "9", "--name", "Renamed"})
	assert.NoError(cmd.Execute())
	assert.Equal(9, gotBody.ID)
	assert.Equal("Renamed", gotBody.Name)
	assert.Equal("https://example.com", gotBody.URL)
	assert.Equal("keep me", gotBody.Description)
	assert.Equal(2, gotBody.Collection.ID)
	assert.Equal(1, gotBody.Collection.OwnerID)
	assert.Len(gotBody.Tags, 1)
	assert.Equal("go", gotBody.Tags[0].Name)
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
