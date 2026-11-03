package links

import (
	"fmt"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestBulkDelete_ParsesIDs(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	first := e2etest.Link(t, client, collection.ID)
	second := e2etest.Link(t, client, collection.ID)
	kept := e2etest.Link(t, client, collection.ID)

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", fmt.Sprintf("%d,%d", first.ID, second.ID)})
	assert.NoError(cmd.Execute())

	links, err := client.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	assert.Equal([]int{kept.ID}, e2etest.LinkIDs(links))
}

func TestBulkDelete_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", "4"})
	assert.Error(cmd.Execute())
}
