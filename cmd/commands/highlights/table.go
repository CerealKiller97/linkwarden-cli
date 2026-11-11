package highlights

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderHighlightTable(app *app.App, highlight linkwarden.Highlight) {
	tables.Highlights(app, "Highlights", []linkwarden.Highlight{highlight}).JSON(highlight).Render()
}
