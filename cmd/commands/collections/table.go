package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderCollectionsTable(collections []linkwarden.Collection) {
	t := utils.NewTable("Collections")
	t.Columns(
		utils.Column{Header: "ID"},
		utils.Column{Header: "Name", Flex: true},
		utils.Column{Header: "Description", Flex: true, HideWhenEmpty: true},
		utils.Column{Header: "Color"},
	)

	for _, c := range collections {
		t.Row(c.ID, c.Name, c.Description, c.Color)
	}

	t.Render()
}
