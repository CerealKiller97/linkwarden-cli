package links

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestDelete_UsesIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	link := e2etest.Link(t, client, collection.ID)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", strconv.Itoa(link.ID)})
	assert.NoError(cmd.Execute())

	links, err := client.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	assert.Empty(links)
}

func TestDelete_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", "999999999"})
	assert.Error(cmd.Execute())
}
