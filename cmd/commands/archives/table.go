package archives

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderLinkTable(link linkwarden.Link) {
	t := utils.NewTable("Link")
	t.AppendHeader(table.Row{"ID", "Name", "URL", "Collection"})
	t.AppendRow(table.Row{link.ID, link.Name, link.URL, link.Collection.Name})
	t.AppendSeparator()
	t.Render()
}
