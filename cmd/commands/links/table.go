package links

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(app *app.App, links []linkwarden.Link) {
	tables.Links(app, "Links", links).Render()
}

// renderLinkTable prints one link, as a JSON object with --json.
func renderLinkTable(app *app.App, link *linkwarden.Link) {
	tables.Links(app, "Links", []linkwarden.Link{*link}).JSON(link).Render()
}

func renderHighlightsTable(app *app.App, highlights []linkwarden.Highlight) {
	tables.Highlights(app, "Highlights", highlights).Render()
}
