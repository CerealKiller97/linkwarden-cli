package search

import (
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderSearchResultsTable(links []linkwarden.Link) {
	t := utils.NewTable("Search Results")
	t.AppendHeader(table.Row{"ID", "Name", "URL", "Collection", "Tags"})

	for _, l := range links {
		tagNames := make([]string, len(l.Tags))
		for i, tag := range l.Tags {
			tagNames[i] = tag.Name
		}

		t.AppendRow(table.Row{l.ID, l.Name, l.URL, l.Collection.Name, strings.Join(tagNames, ", ")})
	}

	t.AppendSeparator()
	t.Render()
}
