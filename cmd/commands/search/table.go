package search

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderSearchResultsTable(app *app.App, links []linkwarden.Link) {
	tables.Links(app, "Search Results", links).Render()
}
