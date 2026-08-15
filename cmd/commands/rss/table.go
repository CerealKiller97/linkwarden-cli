package rss

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderRSSTable(subs []linkwarden.RSSSubscription) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("RSS Subscriptions")
	t.AppendHeader(table.Row{"ID", "Name", "URL", "Collection ID"})

	for _, sub := range subs {
		t.AppendRow(table.Row{sub.ID, sub.Name, sub.URL, sub.CollectionID})
	}

	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
	t.Render()
}
