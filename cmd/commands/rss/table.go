package rss

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderRSSTable(subs []linkwarden.RSSSubscription) {
	t := utils.NewTable("RSS Subscriptions")
	t.AppendHeader(table.Row{"ID", "Name", "URL", "Collection ID"})

	for _, sub := range subs {
		t.AppendRow(table.Row{sub.ID, sub.Name, sub.URL, sub.CollectionID})
	}

	t.AppendSeparator()
	t.Render()
}
