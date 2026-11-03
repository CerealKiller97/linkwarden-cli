package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderTagsTable(tags []linkwarden.TagInfo) {
	t := utils.NewTable("Tags")
	t.AppendHeader(table.Row{"ID", "Name"})

	for _, tag := range tags {
		t.AppendRow(table.Row{tag.ID, tag.Name})
	}

	t.AppendSeparator()
	t.Render()
}
