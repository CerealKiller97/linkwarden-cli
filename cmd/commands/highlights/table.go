package highlights

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderHighlightTable(highlight linkwarden.Highlight) {
	t := utils.NewTable("Highlight")
	t.Header("ID", "Link ID", "Text", "Color", "Comment")
	t.Row(highlight.ID, highlight.LinkID, highlight.Text, highlight.Color, highlight.Comment)
	t.Render()
}
