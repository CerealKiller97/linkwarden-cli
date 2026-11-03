package collections

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func collectionsPickerFetch(ctx context.Context, client *linkwarden.Client) func() ([]utils.PickerItem, error) {
	return func() ([]utils.PickerItem, error) {
		collections, err := client.Collections.List(ctx)
		if err != nil {
			return nil, err
		}

		items := make([]utils.PickerItem, len(collections))
		for i, c := range collections {
			items[i] = utils.PickerItem{Label: c.Name, Value: c.ID}
		}

		return items, nil
	}
}
