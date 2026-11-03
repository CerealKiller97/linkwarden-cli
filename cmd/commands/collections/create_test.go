package collections

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
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

	cmd := Create()
	cmd.SetArgs([]string{"--name", "New Collection"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, gotBody.ParentID)
}

// newUnorganizedServer serves an Unorganized collection with the given ID, counting
// list calls and recording the parent ID of each created collection.
func newUnorganizedServer(t *testing.T, unorganizedID int, listCalls *int, gotParentID *int) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			*listCalls++
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"response": []map[string]any{{"id": unorganizedID, "name": "Unorganized"}},
			})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/collections":
			var body linkwarden.CreateCollectionRequest
			_ = json.NewDecoder(r.Body).Decode(&body)
			*gotParentID = body.ParentID
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":10,"name":"New Collection"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	utils.SetVars("test-token", server.URL)
}

func TestCreate_CachesUnorganizedIDInXDGStateHome(t *testing.T) {
	assert := require.New(t)

	var listCalls, gotParentID int
	newUnorganizedServer(t, 7, &listCalls, &gotParentID)

	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)

	cmd := Create()
	cmd.SetArgs([]string{"--name", "First"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, listCalls)
	assert.Equal(7, gotParentID)

	cached, err := os.ReadFile(filepath.Join(stateHome, "linkwarden", "unorganized"))
	assert.NoError(err)
	assert.Equal("7", string(cached))

	// The second run reads the cached ID back instead of listing collections again.
	cmd = Create()
	cmd.SetArgs([]string{"--name", "Second"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, listCalls)
	assert.Equal(7, gotParentID)
}

func TestCreate_EmptyXDGStateHomeFallsBackToHome(t *testing.T) {
	assert := require.New(t)

	var listCalls, gotParentID int
	newUnorganizedServer(t, 3, &listCalls, &gotParentID)

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // os.UserHomeDir on Windows
	t.Setenv("XDG_STATE_HOME", "")

	cmd := Create()
	cmd.SetArgs([]string{"--name", "New Collection"})
	assert.NoError(cmd.Execute())

	cached, err := os.ReadFile(filepath.Join(home, ".local", "state", "linkwarden", "unorganized"))
	assert.NoError(err)
	assert.Equal("3", string(cached))
}
