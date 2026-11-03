package highlights

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderHighlightTable(highlight linkwarden.Highlight) {
	t := utils.NewTable("Highlight")
	t.AppendHeader(table.Row{"ID", "Link ID", "Text", "Color", "Comment"})
	t.AppendRow(table.Row{highlight.ID, highlight.LinkID, highlight.Text, highlight.Color, highlight.Comment})
	t.AppendSeparator()
	t.Render()
}
