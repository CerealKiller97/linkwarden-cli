package links

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestUpdate_AppliesFlags(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	link := e2etest.Link(t, client, e2etest.Collection(t, client).ID)

	cmd := Update()
	cmd.SetArgs([]string{"--id", strconv.Itoa(link.ID), "--name", "Renamed", "--description", "updated"})
	assert.NoError(cmd.Execute())

	got, err := client.Links.Get(t.Context(), link.ID)
	assert.NoError(err)
	assert.Equal("Renamed", got.Name)
	assert.Equal("updated", got.Description)
	assert.Equal(link.URL, got.URL)
	assert.Equal(link.CollectionID, got.Collection.ID)
}

func TestUpdate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Update()
	cmd.SetArgs([]string{"--id", "999999999", "--name", "Renamed"})
	assert.Error(cmd.Execute())
}
