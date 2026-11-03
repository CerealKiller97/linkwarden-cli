package links

import (
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderLinksTable(links []linkwarden.Link) {
	t := utils.NewTable("Links")
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

func renderHighlightsTable(highlights []linkwarden.Highlight) {
	t := utils.NewTable("Highlights")
	t.AppendHeader(table.Row{"ID", "Text", "Color", "Comment"})

	for _, h := range highlights {
		t.AppendRow(table.Row{h.ID, h.Text, h.Color, h.Comment})
	}

	t.AppendSeparator()
	t.Render()
}
