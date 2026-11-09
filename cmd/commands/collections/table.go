package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

// renderCollectionsTable prints collections; --json prints jsonValue, the whole list or the one
// record a get or create returned.
func renderCollectionsTable(collections []linkwarden.Collection, jsonValue any) {
	t := utils.NewTable("Collections")
	t.Header("ID", "Name", "Description", "Color")

	for _, c := range collections {
		t.Row(c.ID, c.Name, c.Description, utils.ColorSwatch(c.Color))
	}

	t.JSON(jsonValue)
	t.Render()
}
