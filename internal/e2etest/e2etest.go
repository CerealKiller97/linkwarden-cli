// Package e2etest points command tests at a real Linkwarden instance and creates the
// data they act on. Only test files import it.
//
// The instance is set through LINKWARDEN_URL and LINKWARDEN_TOKEN. `scripts/dev-linkwarden.sh up`
// starts a throwaway one locally, and CI starts one per run. Without them, tests are skipped
// locally but fail in CI, so a broken workflow can't pass by skipping everything.
//
// Every helper gives its data a unique name and deletes it when the test ends, so tests can
// share one instance and run in any order.
package e2etest

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

// Setup points the CLI at the test instance and returns an SDK client for arranging and
// checking data around the command under test.
func Setup(t *testing.T) *linkwarden.Client {
	t.Helper()

	baseURL, token := os.Getenv("LINKWARDEN_URL"), os.Getenv("LINKWARDEN_TOKEN")
	if baseURL == "" || token == "" {
		if os.Getenv("CI") != "" {
			t.Fatal("LINKWARDEN_URL and LINKWARDEN_TOKEN must be set in CI")
		}
		t.Skip("LINKWARDEN_URL / LINKWARDEN_TOKEN not set; start an instance with scripts/dev-linkwarden.sh up")
	}

	utils.SetVars(token, baseURL)

	return linkwarden.New(baseURL, token, linkwarden.WithLogger(slog.New(slog.DiscardHandler)))
}

// UseInvalidToken points the CLI at the test instance with a token it rejects, so every
// authenticated request fails with an API error.
func UseInvalidToken(t *testing.T) {
	t.Helper()

	Setup(t)
	utils.SetVars("invalid-token", os.Getenv("LINKWARDEN_URL"))
}

// UniqueName returns prefix followed by a random suffix.
func UniqueName(prefix string) string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)

	return prefix + "-" + hex.EncodeToString(b)
}

// UniqueURL returns a URL that no other link on the instance uses.
func UniqueURL() string {
	return "https://example.com/" + UniqueName("e2e")
}

func Collection(t *testing.T, c *linkwarden.Client) *linkwarden.Collection {
	t.Helper()

	collection, err := c.Collections.Create(t.Context(), linkwarden.CreateCollectionRequest{Name: UniqueName("collection")})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Collections.Delete(t.Context(), collection.ID) })

	return collection
}

func Link(t *testing.T, c *linkwarden.Client, collectionID int) *linkwarden.Link {
	t.Helper()

	link, err := c.Links.Create(t.Context(), linkwarden.CreateLinkRequest{
		Name:       UniqueName("link"),
		URL:        UniqueURL(),
		Collection: linkwarden.Collection{ID: collectionID},
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Links.Delete(t.Context(), link.ID) })

	return link
}

func Tag(t *testing.T, c *linkwarden.Client) *linkwarden.TagInfo {
	t.Helper()

	tags, err := c.Tags.CreateOrUpdateBulk(t.Context(), []linkwarden.CreateOrUpdateTagInput{{Label: UniqueName("tag")}})
	require.NoError(t, err)
	require.Len(t, tags, 1)
	t.Cleanup(func() { _ = c.Tags.Delete(t.Context(), tags[0].ID) })

	return &tags[0]
}

// CleanupTagsNamed deletes, when the test ends, the tags a command creates by name.
func CleanupTagsNamed(t *testing.T, c *linkwarden.Client, names ...string) {
	t.Helper()

	t.Cleanup(func() {
		for _, name := range names {
			tags, err := c.Tags.List(t.Context(), linkwarden.ListTagsOptions{Search: name})
			if err != nil {
				continue
			}
			for _, tag := range tags {
				if tag.Name == name {
					_ = c.Tags.Delete(t.Context(), tag.ID)
				}
			}
		}
	})
}

func Highlight(t *testing.T, c *linkwarden.Client, linkID int) *linkwarden.Highlight {
	t.Helper()

	highlight, err := c.Highlights.Create(t.Context(), linkwarden.CreateHighlightRequest{
		Color:     "yellow",
		EndOffset: 5,
		Text:      "Hello",
		LinkID:    linkID,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.Highlights.Delete(t.Context(), highlight.ID) })

	return highlight
}

// FeedURL is a real feed for RSS subscriptions; the instance fetches it in the background.
const FeedURL = "https://go.dev/blog/feed.atom"

func Subscription(t *testing.T, c *linkwarden.Client, collectionID int) *linkwarden.RSSSubscription {
	t.Helper()

	sub, err := c.RSS.Create(t.Context(), linkwarden.CreateRSSRequest{
		Name:         UniqueName("feed"),
		URL:          FeedURL,
		CollectionID: collectionID,
	})
	require.NoError(t, err)
	t.Cleanup(func() { _ = c.RSS.Delete(t.Context(), sub.ID) })

	return sub
}

// PNG returns a small, valid PNG image.
func PNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for x := range 4 {
		for y := range 4 {
			img.Set(x, y, color.RGBA{R: 245, G: 158, B: 11, A: 255})
		}
	}

	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))

	return buf.Bytes()
}

// WriteFile writes contents to a file named name in a temporary directory and returns its path.
func WriteFile(t *testing.T, name string, contents []byte) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, contents, 0o600))

	return path
}

func LinkIDs(links []linkwarden.Link) []int {
	ids := make([]int, len(links))
	for i, l := range links {
		ids[i] = l.ID
	}

	return ids
}
