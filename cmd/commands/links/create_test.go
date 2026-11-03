package links

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestCreate_UploadsLinksFromFile(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	first, second := e2etest.UniqueURL(), e2etest.UniqueURL()
	file := e2etest.WriteFile(t, "links.json", fmt.Appendf(nil,
		`[{"url":%q,"title":"First"},{"url":%q,"title":"Second"}]`, first, second))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", strconv.Itoa(collection.ID)})
	assert.NoError(cmd.Execute())

	links, err := client.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	t.Cleanup(func() { _ = client.Links.BulkDelete(t.Context(), e2etest.LinkIDs(links)) })

	var urls []string
	for _, l := range links {
		urls = append(urls, l.URL)
	}
	assert.ElementsMatch([]string{first, second}, urls)
}

func TestCreate_EmptyFileErrors(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	file := e2etest.WriteFile(t, "links.json", []byte(`[]`))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", strconv.Itoa(collection.ID)})
	assert.ErrorIs(cmd.Execute(), errFileEmpty)
}

func TestCreate_UnknownCollectionErrorsWithoutCreatingLinks(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	url := e2etest.UniqueURL()
	file := e2etest.WriteFile(t, "links.json", fmt.Appendf(nil, `[{"url":%q,"title":"Example"}]`, url))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "999999999"})
	assert.ErrorIs(cmd.Execute(), errCollectionNotFound)

	links, err := client.Links.List(t.Context(), linkwarden.ListLinksOptions{SearchQueryString: url})
	assert.NoError(err)
	assert.Empty(links)
}
