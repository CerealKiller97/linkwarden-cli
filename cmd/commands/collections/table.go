package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderCollectionsTable(collections []linkwarden.Collection) {
	t := utils.NewTable("Collections")
	t.AppendHeader(table.Row{"ID", "Name", "Description", "Color"})

	for _, c := range collections {
		t.AppendRow(table.Row{c.ID, c.Name, c.Description, c.Color})
	}

	t.AppendSeparator()
	t.Render()
}
