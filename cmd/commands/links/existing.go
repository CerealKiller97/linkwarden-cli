package links

import (
	"context"
	"net/url"
	"strings"

	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// existingURLs returns the normalized URLs of every link already in
// collectionID. The links endpoint is cursor-paginated (the cursor is the last
// link ID seen), so it keeps fetching until a page comes back empty.
func existingURLs(ctx context.Context, client *linkwarden.Client, collectionID int) (map[string]bool, error) {
	seen := make(map[string]bool)
	cursor := 0

	for {
		page, err := client.Links.List(ctx, linkwarden.ListLinksOptions{CollectionID: collectionID, Cursor: cursor})
		if err != nil {
			return nil, err
		}
		if len(page) == 0 {
			return seen, nil
		}

		for _, l := range page {
			seen[normalizeURL(l.URL)] = true
		}

		// Guard against a server that ignores the cursor and repeats a page.
		next := page[len(page)-1].ID
		if next == cursor {
			return seen, nil
		}
		cursor = next
	}
}

// normalizeURL makes URLs that point at the same page compare equal:
// surrounding whitespace and trailing slashes are dropped and the scheme and
// host are lowercased. The rest (path case, query, fragment) is kept as is.
func normalizeURL(raw string) string {
	s := strings.TrimSpace(raw)

	u, err := url.Parse(s)
	if err != nil || u.Host == "" {
		return strings.TrimRight(s, "/")
	}

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)
	u.Path = strings.TrimRight(u.Path, "/")
	u.RawPath = ""

	return u.String()
}
