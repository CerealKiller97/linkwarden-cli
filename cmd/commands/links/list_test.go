package links

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestList_RendersLinks(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	e2etest.Link(t, client, collection.ID)

	cmd := List()
	cmd.SetArgs([]string{"--collection-id", strconv.Itoa(collection.ID)})
	assert.NoError(cmd.Execute())
}

func TestList_PinnedOnly(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	e2etest.Link(t, client, collection.ID)

	cmd := List()
	cmd.SetArgs([]string{"--collection-id", strconv.Itoa(collection.ID), "--pinned-only"})
	assert.NoError(cmd.Execute())
}

func TestList_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := List()
	cmd.SetArgs([]string{})
	assert.Error(cmd.Execute())
}
