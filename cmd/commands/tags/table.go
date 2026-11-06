package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderTagsTable(tags []linkwarden.TagInfo) {
	t := utils.NewTable("Tags")
	t.Header("ID", "Name")

	for _, tag := range tags {
		t.Row(tag.ID, tag.Name)
	}

	t.Render()
}
