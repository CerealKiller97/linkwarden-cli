package rss

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// colRSSURL is the URL column of the RSS table.
const colRSSURL = 2

func renderRSSTable(subs []linkwarden.RSSSubscription) {
	t := utils.NewTable("RSS Subscriptions")
	t.Header("ID", "Name", "URL", "Collection ID")

	for _, sub := range subs {
		t.Row(sub.ID, sub.Name, sub.URL, sub.CollectionID)
	}
	t.LinkColumn(colRSSURL)

	t.Render()
}
