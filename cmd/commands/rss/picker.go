package rss

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func rssPickerFetch(ctx context.Context, client *linkwarden.Client) func() ([]utils.PickerItem, error) {
	return func() ([]utils.PickerItem, error) {
		subs, err := client.RSS.List(ctx)
		if err != nil {
			return nil, err
		}

		items := make([]utils.PickerItem, len(subs))
		for i, sub := range subs {
			items[i] = utils.PickerItem{Label: sub.Name, Value: sub.ID}
		}

		return items, nil
	}
}
