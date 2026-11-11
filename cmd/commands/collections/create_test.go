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

func TestCreateUsesParentIDFlagWithoutListing(t *testing.T) {
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

	cmd := Create(nil)
	cmd.SetArgs([]string{"--name", "New Collection", "--parent-id", "5"})
	assert.NoError(cmd.Execute())
	assert.False(listCalled)
	assert.Equal("New Collection", gotBody.Name)
	assert.Equal(5, gotBody.ParentID)
}

func TestCreateLooksUpUnorganizedWhenParentIDOmitted(t *testing.T) {
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

	cmd := Create(nil)
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

func TestCreateCachesUnorganizedIDPerInstance(t *testing.T) {
	assert := require.New(t)

	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)

	createIn := func(name string) {
		cmd := Create(nil)
		cmd.SetArgs([]string{"--name", name})
		assert.NoError(cmd.Execute())
	}

	// Two instances whose Unorganized collections have different IDs.
	var localLists, localParent, cloudLists, cloudParent int
	newUnorganizedServer(t, 8, &localLists, &localParent)
	localURL := utils.GetInstanceURL()
	newUnorganizedServer(t, 38899, &cloudLists, &cloudParent)
	cloudURL := utils.GetInstanceURL()

	utils.SetVars("test-token", localURL)
	createIn("First")
	createIn("Second")
	assert.Equal(1, localLists, "the second run reads the cached ID")
	assert.Equal(8, localParent)

	// Switching instances must not reuse the other instance's ID.
	utils.SetVars("test-token", cloudURL)
	createIn("Third")
	assert.Equal(1, cloudLists)
	assert.Equal(38899, cloudParent)

	for url, want := range map[string]string{localURL: "8", cloudURL: "38899"} {
		file, err := utils.InstanceStateFile("unorganized", url)
		assert.NoError(err)
		assert.Equal(filepath.Join(stateHome, "linkwarden", "unorganized"), filepath.Dir(file))

		cached, err := os.ReadFile(file)
		assert.NoError(err)
		assert.Equal(want, string(cached))
	}
}

func TestCreateReplacesOldSingleCacheFile(t *testing.T) {
	assert := require.New(t)

	// Older versions kept one ID for every instance in this file. Using it sent
	// another instance's collection ID as the parent, which Linkwarden rejects.
	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)
	legacy := filepath.Join(stateHome, "linkwarden", "unorganized")
	assert.NoError(os.MkdirAll(filepath.Dir(legacy), 0o750))
	assert.NoError(os.WriteFile(legacy, []byte("8"), 0o600))

	var listCalls, gotParentID int
	newUnorganizedServer(t, 38899, &listCalls, &gotParentID)

	cmd := Create(nil)
	cmd.SetArgs([]string{"--name", "New Collection"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, listCalls, "the old file is ignored, not trusted")
	assert.Equal(38899, gotParentID)

	info, err := os.Stat(legacy)
	assert.NoError(err)
	assert.True(info.IsDir(), "the old file is replaced by the per-instance directory")
}

func TestUnorganizedParentID(t *testing.T) {
	assert := require.New(t)

	collections := []linkwarden.Collection{
		{ID: 1, Name: "Unorganized"},
		{ID: 2, Name: "Reading"},
	}

	id, err := unorganizedParentID(collections, "New")
	assert.NoError(err)
	assert.Equal(1, id)

	_, err = unorganizedParentID(collections, "Reading")
	assert.ErrorContains(err, `collection "Reading" already exists`, "listed after Unorganized, still caught")

	id, err = unorganizedParentID([]linkwarden.Collection{{ID: 2, Name: "Reading"}}, "New")
	assert.NoError(err)
	assert.Zero(id, "no Unorganized collection means a top-level collection")
}
