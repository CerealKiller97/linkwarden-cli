package highlights

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderHighlightTable(highlight linkwarden.Highlight) {
	t := utils.NewTable("Highlight")
	t.Columns(
		utils.Column{Header: "ID"},
		utils.Column{Header: "Link ID"},
		utils.Column{Header: "Text", Flex: true},
		utils.Column{Header: "Color"},
		utils.Column{Header: "Comment", Flex: true, HideWhenEmpty: true},
	)
	t.Row(highlight.ID, highlight.LinkID, highlight.Text, highlight.Color, highlight.Comment)
	t.Render()
}
