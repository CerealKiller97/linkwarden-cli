package rss

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestCreate_SubscribesToFeed(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	name := e2etest.UniqueName("feed")

	cmd := Create()
	cmd.SetArgs([]string{"--name", name, "--url", e2etest.FeedURL, "--collection-id", strconv.Itoa(collection.ID)})
	assert.NoError(cmd.Execute())

	subs, err := client.RSS.List(t.Context())
	assert.NoError(err)

	var found bool
	for _, sub := range subs {
		if sub.Name == name {
			found = true
			t.Cleanup(func() { _ = client.RSS.Delete(t.Context(), sub.ID) })
			assert.Equal(e2etest.FeedURL, sub.URL)
			assert.Equal(collection.ID, sub.CollectionID)
		}
	}
	assert.True(found, "new subscription should be listed")
}

func TestCreate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := Create()
	cmd.SetArgs([]string{"--name", "News", "--url", e2etest.FeedURL})
	assert.Error(cmd.Execute())
}
