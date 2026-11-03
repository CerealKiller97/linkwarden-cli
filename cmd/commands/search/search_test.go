package search

import (
	"strings"
	"testing"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestExecute_FindsIndexedLink(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	// One word with no separators, so the search can't match other links on a shared prefix.
	word := strings.ReplaceAll(e2etest.UniqueName("searchable"), "-", "")
	link, err := client.Links.Create(t.Context(), linkwarden.CreateLinkRequest{
		Name:       word,
		URL:        e2etest.UniqueURL(),
		Collection: linkwarden.Collection{ID: collection.ID},
	})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Links.Delete(t.Context(), link.ID) })

	// Wait until the link reaches the search index, then run the command against it.
	assert.Eventually(func() bool {
		links, err := client.Search.Links(t.Context(), linkwarden.SearchOptions{SearchQueryString: word})
		return err == nil && len(links) > 0
	}, time.Minute, time.Second)

	cmd := Command()
	cmd.SetArgs([]string{"--query", word, "--sort", "0"})
	assert.NoError(cmd.Execute())
}

func TestExecute_NoMatches(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Command()
	cmd.SetArgs([]string{"--query", e2etest.UniqueName("nothing")})
	assert.NoError(cmd.Execute())
}

func TestExecute_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := Command()
	cmd.SetArgs([]string{"--query", "go"})
	assert.Error(cmd.Execute())
}
