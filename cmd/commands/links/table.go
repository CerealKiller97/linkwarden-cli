package links

import (
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(links []linkwarden.Link) {
	t := utils.NewTable("Links")
	t.Header("ID", "Name", "URL", "Collection", "Tags")

	for _, l := range links {
		tagNames := make([]string, len(l.Tags))
		for i, tag := range l.Tags {
			tagNames[i] = tag.Name
		}

		t.Row(l.ID, l.Name, l.URL, l.Collection.Name, strings.Join(tagNames, ", "))
	}

	t.Render()
}

func renderHighlightsTable(highlights []linkwarden.Highlight) {
	t := utils.NewTable("Highlights")
	t.Header("ID", "Text", "Color", "Comment")

	for _, h := range highlights {
		t.Row(h.ID, h.Text, h.Color, h.Comment)
	}

	t.Render()
}
