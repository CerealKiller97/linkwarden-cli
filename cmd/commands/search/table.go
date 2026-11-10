package search

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderSearchResultsTable(ctx context.Context, links []linkwarden.Link) {
	tables.Links(ctx, "Search Results", links).Render()
}
