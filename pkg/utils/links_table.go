package utils

import (
	"strings"

	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// NewLinksTable builds the table shared by every command that lists links
// (links list, search, dashboard): NAME and URL shrink to fit the terminal,
// and TAGS is left out when no link has any.
func NewLinksTable(title string, links []linkwarden.Link) *Table {
	t := NewTable(title)
	t.Columns(
		Column{Header: "ID"},
		Column{Header: "Name", Flex: true, MinWidth: 20},
		Column{Header: "URL", Flex: true, MinWidth: 20},
		Column{Header: "Collection", Flex: true},
		Column{Header: "Tags", Flex: true, HideWhenEmpty: true},
	)

	for _, l := range links {
		tagNames := make([]string, len(l.Tags))
		for i, tag := range l.Tags {
			tagNames[i] = tag.Name
		}

		t.Row(l.ID, l.Name, l.URL, l.Collection.Name, strings.Join(tagNames, ", "))
	}

	return t
}

// RenderLinksTable renders NewLinksTable.
func RenderLinksTable(title string, links []linkwarden.Link) {
	NewLinksTable(title, links).Render()
}
