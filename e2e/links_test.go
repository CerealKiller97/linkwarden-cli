//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/stretchr/testify/require"
)

func TestLinks_CreateFromFile(t *testing.T) {
	assert := require.New(t)

	collection := newCollection(t)
	marker := unique("link")
	items, err := json.Marshal([]map[string]string{
		{"url": "https://example.com/" + marker + "/1", "title": marker + " one"},
		{"url": "https://example.com/" + marker + "/2", "title": marker + " two"},
		{"url": "not a url", "title": marker + " bad"},
	})
	assert.NoError(err)
	file := writeFile(t, "links.json", string(items))

	out := mustCLI(t, "links", "create", "--file", file, "--collection", itoa(collection.ID))
	assert.Regexp(`│\s+3\s+│\s+2\s+│\s+1\s+│`, out, "expected total=3 success=2 failed=1")

	links, err := api.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	assert.Len(links, 2)

	out = mustCLI(t, "links", "list", "--collection-id", itoa(collection.ID))
	assert.Contains(out, marker+" one")
	assert.Contains(out, marker+" two")
}

func TestLinks_GetAndDelete(t *testing.T) {
	assert := require.New(t)

	link := newLink(t, newCollection(t), linkwarden.CreateLinkRequest{})

	assert.Contains(mustCLI(t, "links", "get", "--id", itoa(link.ID)), link.Name)

	mustCLI(t, "links", "delete", "--id", itoa(link.ID))
	_, err := api.Links.Get(t.Context(), link.ID)
	assert.Error(err, "link still exists after delete")
}

func TestLinks_UpdateKeepsUnchangedFields(t *testing.T) {
	assert := require.New(t)

	collection := newCollection(t)
	tagName := unique("tag")
	link := newLink(t, collection, linkwarden.CreateLinkRequest{
		Description: "original description",
		Tags:        []linkwarden.Tag{{Name: tagName}},
	})
	cleanup(t, func(ctx context.Context) error {
		tags, err := api.Tags.List(ctx, linkwarden.ListTagsOptions{Search: tagName})
		for _, tag := range tags {
			_ = api.Tags.Delete(ctx, tag.ID)
		}
		return err
	})

	mustCLI(t, "links", "update", "--id", itoa(link.ID), "--name", link.Name+"-renamed", "--color", "red")

	got, err := api.Links.Get(t.Context(), link.ID)
	assert.NoError(err)
	assert.Equal(link.Name+"-renamed", got.Name)
	assert.Equal("red", got.Color)
	assert.Equal("original description", got.Description)
	assert.Equal(link.URL, got.URL)
	assert.Equal(collection.ID, got.Collection.ID)
	assert.Equal([]string{tagName}, tagNames(got.Tags))
}

func TestLinks_BulkUpdateAndBulkDelete(t *testing.T) {
	assert := require.New(t)

	collection := newCollection(t)
	first := newLink(t, collection, linkwarden.CreateLinkRequest{})
	second := newLink(t, collection, linkwarden.CreateLinkRequest{})
	ids := fmt.Sprintf("%d,%d", first.ID, second.ID)

	tagName := unique("bulk")
	mustCLI(t, "links", "bulk-update", "--ids", ids, "--tag", tagName)
	tag := findTag(t, tagName)
	cleanup(t, func(ctx context.Context) error { return api.Tags.Delete(ctx, tag.ID) })

	for _, link := range []linkwarden.Link{first, second} {
		got, err := api.Links.Get(t.Context(), link.ID)
		assert.NoError(err)
		assert.Contains(tagNames(got.Tags), tagName)
	}

	out := mustCLI(t, "links", "list", "--tag-id", itoa(tag.ID))
	assert.Contains(out, first.Name)
	assert.Contains(out, second.Name)

	mustCLI(t, "links", "bulk-delete", "--ids", ids)
	links, err := api.Links.List(t.Context(), linkwarden.ListLinksOptions{CollectionID: collection.ID})
	assert.NoError(err)
	assert.Empty(links)
}

func TestSearch(t *testing.T) {
	assert := require.New(t)

	link := newLink(t, newCollection(t), linkwarden.CreateLinkRequest{})

	// Search goes through the server's index, which is updated asynchronously.
	eventually(t, 30*time.Second, "search never found "+link.Name, func() bool {
		out, code := cliWithEnv(t, nil, "search", "--query", link.Name)
		return code == 0 && strings.Contains(out, link.Name)
	})

	// The API answers a search with no hits with "data": [] instead of an object.
	// The query shares no words with the e2e- names, which fuzzy search would match.
	out := mustCLI(t, "search", "--query", "zqxv"+strconv.FormatInt(time.Now().UnixNano(), 36))
	assert.Contains(out, "Search Results")
	assert.NotRegexp(`│\s+\d+\s+│`, out, "expected no results")
}

func TestHighlights(t *testing.T) {
	assert := require.New(t)

	link := newLink(t, newCollection(t), linkwarden.CreateLinkRequest{})
	comment := unique("comment")

	mustCLI(t, "highlights", "create", "--link-id", itoa(link.ID), "--text", "Example",
		"--start-offset", "0", "--end-offset", "7", "--color", "yellow", "--comment", comment)

	highlights, err := api.Links.Highlights(t.Context(), link.ID)
	assert.NoError(err)
	assert.Len(highlights, 1)
	assert.Equal(comment, highlights[0].Comment)

	assert.Contains(mustCLI(t, "links", "get-highlights", "--id", itoa(link.ID)), comment)

	mustCLI(t, "highlights", "delete", "--id", itoa(highlights[0].ID))
	highlights, err = api.Links.Highlights(t.Context(), link.ID)
	assert.NoError(err)
	assert.Empty(highlights)
}

func tagNames(tags []linkwarden.Tag) []string {
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}
	slices.Sort(names)

	return names
}
