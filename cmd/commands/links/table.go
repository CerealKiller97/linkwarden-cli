package links

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(links []linkwarden.Link) {
	utils.RenderLinksTable("Links", links)
}

func renderHighlightsTable(highlights []linkwarden.Highlight) {
	t := utils.NewTable("Highlights")
	t.Columns(
		utils.Column{Header: "ID"},
		utils.Column{Header: "Text", Flex: true},
		utils.Column{Header: "Color"},
		utils.Column{Header: "Comment", Flex: true, HideWhenEmpty: true},
	)

	for _, h := range highlights {
		t.Row(h.ID, h.Text, h.Color, h.Comment)
	}

	t.Render()
}
