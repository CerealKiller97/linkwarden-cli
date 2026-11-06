package rss

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderRSSTable(subs []linkwarden.RSSSubscription) {
	t := utils.NewTable("RSS Subscriptions")
	t.Columns(
		utils.Column{Header: "ID"},
		utils.Column{Header: "Name", Flex: true},
		utils.Column{Header: "URL", Flex: true},
		utils.Column{Header: "Collection ID"},
	)

	for _, sub := range subs {
		t.Row(sub.ID, sub.Name, sub.URL, sub.CollectionID)
	}

	t.Render()
}
