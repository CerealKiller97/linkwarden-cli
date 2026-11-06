package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderTagsTable(tags []linkwarden.TagInfo) {
	t := utils.NewListTable("Tags",
		utils.Column{Header: "ID"},
		utils.Column{Header: "Name", Flex: true},
	)

	for _, tag := range tags {
		t.AppendRow(tag.ID, tag.Name)
	}

	t.Render()
}
