package highlights

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderHighlightTable(highlight linkwarden.Highlight) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Highlight")
	t.AppendHeader(table.Row{"ID", "Link ID", "Text", "Color", "Comment"})
	t.AppendRow(table.Row{highlight.ID, highlight.LinkID, highlight.Text, highlight.Color, highlight.Comment})
	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
	t.Render()
}
