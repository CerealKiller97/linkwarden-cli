package collections

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

// unorganizedID returns the ID of the user's Unorganized collection, creating it if this
// instance doesn't have one yet. It is shared by every test and never deleted.
func unorganizedID(t *testing.T, c *linkwarden.Client) int {
	t.Helper()

	collections, err := c.Collections.List(t.Context())
	require.NoError(t, err)
	for _, collection := range collections {
		if collection.Name == unorganizedCollectionName {
			return collection.ID
		}
	}

	created, err := c.Collections.Create(t.Context(), linkwarden.CreateCollectionRequest{Name: unorganizedCollectionName})
	require.NoError(t, err)

	return created.ID
}

// createdCollection finds the collection a command created under name, and deletes it
// when the test ends.
func createdCollection(t *testing.T, c *linkwarden.Client, name string) *linkwarden.Collection {
	t.Helper()

	collections, err := c.Collections.List(t.Context())
	require.NoError(t, err)
	for i := range collections {
		if collections[i].Name == name {
			t.Cleanup(func() { _ = c.Collections.Delete(t.Context(), collections[i].ID) })
			return &collections[i]
		}
	}

	t.Fatalf("no collection named %q was created", name)

	return nil
}

func TestCreate_UsesParentIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	parent := e2etest.Collection(t, client)
	name := e2etest.UniqueName("child")

	cmd := Create()
	cmd.SetArgs([]string{"--name", name, "--parent-id", strconv.Itoa(parent.ID), "--color", "#f59e0b"})
	assert.NoError(cmd.Execute())

	created := createdCollection(t, client, name)
	assert.NotNil(created.ParentID)
	assert.Equal(parent.ID, *created.ParentID)
	assert.Equal("#f59e0b", created.Color)
}

func TestCreate_LooksUpUnorganizedWhenParentIDOmitted(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	unorganized := unorganizedID(t, client)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	name := e2etest.UniqueName("collection")

	cmd := Create()
	cmd.SetArgs([]string{"--name", name})
	assert.NoError(cmd.Execute())

	created := createdCollection(t, client, name)
	assert.NotNil(created.ParentID)
	assert.Equal(unorganized, *created.ParentID)
}

func TestCreate_CachesUnorganizedIDInXDGStateHome(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	unorganized := unorganizedID(t, client)
	stateHome := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateHome)

	first := e2etest.UniqueName("first")
	cmd := Create()
	cmd.SetArgs([]string{"--name", first})
	assert.NoError(cmd.Execute())
	createdCollection(t, client, first)

	cached, err := os.ReadFile(filepath.Join(stateHome, "linkwarden", "unorganized"))
	assert.NoError(err)
	assert.Equal(strconv.Itoa(unorganized), string(cached))

	// The second run reads the cached ID back and still lands under Unorganized.
	second := e2etest.UniqueName("second")
	cmd = Create()
	cmd.SetArgs([]string{"--name", second})
	assert.NoError(cmd.Execute())

	created := createdCollection(t, client, second)
	assert.NotNil(created.ParentID)
	assert.Equal(unorganized, *created.ParentID)
}

func TestCreate_EmptyXDGStateHomeFallsBackToHome(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	unorganized := unorganizedID(t, client)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // os.UserHomeDir on Windows
	t.Setenv("XDG_STATE_HOME", "")
	name := e2etest.UniqueName("collection")

	cmd := Create()
	cmd.SetArgs([]string{"--name", name})
	assert.NoError(cmd.Execute())
	createdCollection(t, client, name)

	cached, err := os.ReadFile(filepath.Join(home, ".local", "state", "linkwarden", "unorganized"))
	assert.NoError(err)
	assert.Equal(strconv.Itoa(unorganized), string(cached))
}
