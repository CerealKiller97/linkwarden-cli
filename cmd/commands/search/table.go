package search

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderSearchResultsTable(links []linkwarden.Link) {
	tables.Links("Search Results", links).Render()
}
