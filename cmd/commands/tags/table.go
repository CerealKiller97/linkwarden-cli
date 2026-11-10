package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// renderTagsTable prints tags; --json prints jsonValue, the whole list or the one
// record a get or create returned.
func renderTagsTable(app *app.App, tags []linkwarden.TagInfo, jsonValue any) {
	t := utils.NewTable(app, "Tags")
	t.Header("ID", "Name")

	for _, tag := range tags {
		t.Row(tag.ID, tag.Name)
	}

	t.JSON(jsonValue)
	t.Render()
}
