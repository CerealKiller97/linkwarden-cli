package collections

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderCollectionsTable(collections []linkwarden.Collection) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Collections")
	t.AppendHeader(table.Row{"ID", "Name", "Description", "Color"})

	for _, c := range collections {
		t.AppendRow(table.Row{c.ID, c.Name, c.Description, c.Color})
	}

	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
	t.Render()
}
