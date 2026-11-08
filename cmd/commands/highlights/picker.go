package highlights

import (
	"context"
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

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
