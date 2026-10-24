//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/stretchr/testify/require"
)

func TestCollections_Lifecycle(t *testing.T) {
	assert := require.New(t)

	name := unique("collection")
	out := mustCLI(t, "collections", "create", "--name", name, "--description", "original", "--parent-id", "0")
	id := loggedID(t, out)
	cleanup(t, func(ctx context.Context) error { return api.Collections.Delete(ctx, id) })

	assert.Contains(mustCLI(t, "collections", "get", "--id", itoa(id)), name)
	assert.Contains(mustCLI(t, "collections", "list"), name)

	out = failCLI(t, "collections", "create", "--name", name, "--parent-id", "0")
	assert.Contains(out, "already exists")

	mustCLI(t, "collections", "delete", "--id", itoa(id))
	_, err := api.Collections.Get(t.Context(), id)
	assert.True(linkwarden.IsNotFound(err), "collection still exists after delete: %v", err)
}

func TestCollections_UpdateKeepsUnchangedFields(t *testing.T) {
	assert := require.New(t)

	parent := newCollection(t)
	collection, err := api.Collections.Create(t.Context(), linkwarden.CreateCollectionRequest{
		Name:        unique("collection"),
		Description: "original",
		Color:       "#123456",
	})
	assert.NoError(err)
	cleanup(t, func(ctx context.Context) error { return api.Collections.Delete(ctx, collection.ID) })
	id := itoa(collection.ID)

	mustCLI(t, "collections", "update", "--id", id, "--name", collection.Name+"-renamed")
	got, err := api.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.Equal(collection.Name+"-renamed", got.Name)
	assert.Equal("original", got.Description)
	assert.Equal("#123456", got.Color)

	mustCLI(t, "collections", "update", "--id", id, "--public", "--parent-id", itoa(parent.ID))
	got, err = api.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.True(got.IsPublic)
	assert.NotNil(got.ParentID)
	assert.Equal(parent.ID, *got.ParentID)
	assert.Equal("original", got.Description)

	mustCLI(t, "collections", "update", "--id", id, "--public=false")
	got, err = api.Collections.Get(t.Context(), collection.ID)
	assert.NoError(err)
	assert.False(got.IsPublic)
	assert.Equal(collection.Name+"-renamed", got.Name)
}

func TestCollections_GetMissingFails(t *testing.T) {
	out := failCLI(t, "collections", "get", "--id", "2147483000")
	require.Contains(t, out, "not found")
}
