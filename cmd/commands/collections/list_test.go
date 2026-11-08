package collections

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestList_RendersCollections(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/collections", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"}]}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	cmd := List()
	assert.NoError(cmd.Execute())
}

func TestList_CachesUnorganizedIDForCreate(t *testing.T) {
	assert := require.New(t)

	t.Setenv("XDG_STATE_HOME", t.TempDir())

	var listCalls, gotParentID int
	newUnorganizedServer(t, 38899, &listCalls, &gotParentID)

	assert.NoError(List().Execute())
	assert.Equal(1, listCalls)

	file, err := utils.InstanceStateFile("unorganized", utils.GetInstanceURL())
	assert.NoError(err)
	cached, err := os.ReadFile(file)
	assert.NoError(err)
	assert.Equal("38899", string(cached))

	// create then uses the ID list cached, without listing again.
	cmd := Create()
	cmd.SetArgs([]string{"--name", "New Collection"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, listCalls)
	assert.Equal(38899, gotParentID)
}

func TestList_SucceedsWhenCacheCantBeWritten(t *testing.T) {
	assert := require.New(t)

	// A file where the state directory should be makes the cache unwritable.
	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)
	assert.NoError(os.WriteFile(filepath.Join(stateHome, "linkwarden"), nil, 0o600))

	var listCalls, gotParentID int
	newUnorganizedServer(t, 1, &listCalls, &gotParentID)

	assert.NoError(List().Execute())
}

func TestList_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := List()
	assert.Error(cmd.Execute())
}
