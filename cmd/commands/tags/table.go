package tags

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// renderTagsTable prints tags; --json prints jsonValue, the whole list or the one
// record a get or create returned.
func renderTagsTable(ctx context.Context, tags []linkwarden.TagInfo, jsonValue any) {
	t := utils.NewTable(ctx, "Tags")
	t.Header("ID", "Name")

	for _, tag := range tags {
		t.Row(tag.ID, tag.Name)
	}

	t.JSON(jsonValue)
	t.Render()
}
