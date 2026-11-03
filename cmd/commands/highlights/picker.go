package highlights

import (
	"context"
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
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

// highlightsPickerFetch lists highlights on a single link, since there is no
// list-all-highlights endpoint.
func highlightsPickerFetch(ctx context.Context, client *linkwarden.Client, linkID int) func() ([]utils.PickerItem, error) {
	return func() ([]utils.PickerItem, error) {
		highlights, err := client.Links.Highlights(ctx, linkID)
		if err != nil {
			return nil, err
		}

		items := make([]utils.PickerItem, len(highlights))
		for i, h := range highlights {
			text := h.Text
			if len(text) > 60 {
				text = text[:60] + "…"
			}
			items[i] = utils.PickerItem{Label: fmt.Sprintf("[%d] %s", h.ID, text), Value: h.ID}
		}

		return items, nil
	}
}
