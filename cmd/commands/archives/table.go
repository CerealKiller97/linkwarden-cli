package archives

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderLinkTable(link linkwarden.Link) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Link")
	t.AppendHeader(table.Row{"ID", "Name", "URL", "Collection"})
	t.AppendRow(table.Row{link.ID, link.Name, link.URL, link.Collection.Name})
	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
	t.Render()
}
