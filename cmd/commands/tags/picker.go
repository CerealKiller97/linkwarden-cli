package tags

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func tagsPickerFetch(ctx context.Context, client *linkwarden.Client) func() ([]utils.PickerItem, error) {
	return func() ([]utils.PickerItem, error) {
		tags, err := client.Tags.List(ctx, linkwarden.ListTagsOptions{})
		if err != nil {
			return nil, err
		}

		items := make([]utils.PickerItem, len(tags))
		for i, tag := range tags {
			items[i] = utils.PickerItem{Label: tag.Name, Value: tag.ID}
		}

		return items, nil
	}
}
