// Package tables holds the result tables more than one command prints, so the
// same kind of record always looks the same.
package tables

import (
	"context"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// colLinkURL is the URL column of the Links table.
const colLinkURL = 2

// Links builds the table of links that links, search, dashboard and archives
// print. Call Render on it, after adding a Caption if needed. --json prints
// links; call JSON to print something else, such as the one link of a get.
func Links(ctx context.Context, title string, links []linkwarden.Link) *utils.Table {
	t := utils.NewTable(ctx, title)
	t.Header("ID", "Name", "URL", "Collection", "Tags")

	for _, l := range links {
		t.Row(l.ID, l.Name, l.URL, l.Collection.Name, tagNames(l.Tags))
	}
	t.LinkColumn(colLinkURL)
	t.JSON(links)

	return t
}

// Highlights builds the table of highlights that highlights and links print.
func Highlights(ctx context.Context, title string, highlights []linkwarden.Highlight) *utils.Table {
	t := utils.NewTable(ctx, title)
	t.Header("ID", "Link ID", "Text", "Color", "Comment")

	for _, h := range highlights {
		t.Row(h.ID, h.LinkID, h.Text, h.Color, h.Comment)
	}
	t.JSON(highlights)

	return t
}

func tagNames(tags []linkwarden.Tag) string {
	names := make([]string, len(tags))
	for i, tag := range tags {
		names[i] = tag.Name
	}

	return strings.Join(names, ", ")
}
