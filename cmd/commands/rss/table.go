package rss

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// colRSSURL is the URL column of the RSS table.
const colRSSURL = 2

// renderRSSTable prints subs; --json prints jsonValue, the whole list or the one
// record a get or create returned.
func renderRSSTable(ctx context.Context, subs []linkwarden.RSSSubscription, jsonValue any) {
	t := utils.NewTable(ctx, "RSS Subscriptions")
	t.Header("ID", "Name", "URL", "Collection ID")

	for _, sub := range subs {
		t.Row(sub.ID, sub.Name, sub.URL, sub.CollectionID)
	}
	t.LinkColumn(colRSSURL)

	t.JSON(jsonValue)
	t.Render()
}
