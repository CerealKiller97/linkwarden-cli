//go:build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/stretchr/testify/require"
)

func TestTags_Lifecycle(t *testing.T) {
	assert := require.New(t)

	nameA, nameB := unique("tag-a"), unique("tag-b")
	out := mustCLI(t, "tags", "create", "--label", nameA, "--label", nameB)
	assert.Contains(out, nameA)
	assert.Contains(out, nameB)

	tagA, tagB := findTag(t, nameA), findTag(t, nameB)
	cleanup(t, func(ctx context.Context) error { return api.Tags.BulkDelete(ctx, []int{tagA.ID, tagB.ID}) })

	out = mustCLI(t, "tags", "list", "--search", nameA)
	assert.Contains(out, nameA)
	assert.NotContains(out, nameB)

	assert.Contains(mustCLI(t, "tags", "get", "--id", itoa(tagA.ID)), nameA)

	mustCLI(t, "tags", "update", "--id", itoa(tagA.ID), "--name", nameA+"-renamed")
	got, err := api.Tags.Get(t.Context(), tagA.ID)
	assert.NoError(err)
	assert.Equal(nameA+"-renamed", got.Name)

	mustCLI(t, "tags", "delete", "--id", itoa(tagA.ID))
	_, err = api.Tags.Get(t.Context(), tagA.ID)
	assert.True(linkwarden.IsNotFound(err), "tag still exists after delete: %v", err)

	mustCLI(t, "tags", "bulk-delete", "--ids", itoa(tagB.ID))
	_, err = api.Tags.Get(t.Context(), tagB.ID)
	assert.True(linkwarden.IsNotFound(err), "tag still exists after bulk-delete: %v", err)
}
