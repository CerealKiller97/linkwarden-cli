// Package pickers holds the interactive pickers more than one command uses, so
// picking the same kind of record works the same everywhere.
package pickers

import (
	"context"
	"net/url"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/charmbracelet/x/ansi"
)

// Links fetches the link picker a page at a time, newest first, showing each
// link's name with its shortened URL. Like the API, the cursor is the ID of
// the last link of the previous page.
func Links(ctx context.Context, client *linkwarden.Client) utils.PageFetch {
	return func(cursor int) ([]utils.PickerItem, int, error) {
		links, err := client.Links.List(ctx, linkwarden.ListLinksOptions{Cursor: cursor})
		if err != nil {
			return nil, 0, err
		}
		if len(links) == 0 {
			return nil, 0, nil
		}

		items := make([]utils.PickerItem, len(links))
		for i, l := range links {
			label := strings.TrimSpace(l.Name)
			if label == "" {
				label = l.URL
			}
			items[i] = utils.PickerItem{Label: label, Detail: ShortURL(l.URL), Value: l.ID}
		}

		next := links[len(links)-1].ID
		if next == cursor {
			// A server that ignores the cursor would repeat this page forever.
			next = 0
		}

		return items, next, nil
	}
}

// shortURLWidth caps how much of a URL the picker shows.
const shortURLWidth = 40

// ShortURL is raw as a picker shows it: without the scheme, a leading "www."
// or a trailing slash, and cut to shortURLWidth with an ellipsis.
func ShortURL(raw string) string {
	s := strings.TrimSpace(raw)
	if u, err := url.Parse(s); err == nil && u.Host != "" {
		s = strings.TrimPrefix(u.Host, "www.") + u.EscapedPath()
		if u.RawQuery != "" {
			s += "?" + u.RawQuery
		}
	}

	return ansi.Truncate(strings.TrimRight(s, "/"), shortURLWidth, "…")
}
