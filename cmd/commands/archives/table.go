package archives

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinkTable(app *app.App, link linkwarden.Link) {
	tables.Links(app, "Links", []linkwarden.Link{link}).JSON(link).Render()
}
