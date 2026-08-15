package tags

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderTagsTable(tags []linkwarden.TagInfo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Tags")
	t.AppendHeader(table.Row{"ID", "Name"})

	for _, tag := range tags {
		t.AppendRow(table.Row{tag.ID, tag.Name})
	}

	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
	t.Render()
}
