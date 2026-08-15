package archives

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
)

func linksPickerFetch(ctx context.Context, client *linkwarden.Client) func() ([]utils.PickerItem, error) {
	return func() ([]utils.PickerItem, error) {
		links, err := client.Links.List(ctx, linkwarden.ListLinksOptions{})
		if err != nil {
			return nil, err
		}

		items := make([]utils.PickerItem, len(links))
		for i, l := range links {
			items[i] = utils.PickerItem{Label: l.Name, Value: l.ID}
		}

		return items, nil
	}
}
