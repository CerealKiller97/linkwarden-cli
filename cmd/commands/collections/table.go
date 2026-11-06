package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderCollectionsTable(collections []linkwarden.Collection) {
	t := utils.NewTable("Collections")
	t.Header("ID", "Name", "Description", "Color")

	for _, c := range collections {
		t.Row(c.ID, c.Name, c.Description, c.Color)
	}

	t.Render()
}
