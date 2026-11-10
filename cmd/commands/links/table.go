package links

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(ctx context.Context, links []linkwarden.Link) {
	tables.Links(ctx, "Links", links).Render()
}

// renderLinkTable prints one link, as a JSON object with --json.
func renderLinkTable(ctx context.Context, link *linkwarden.Link) {
	tables.Links(ctx, "Links", []linkwarden.Link{*link}).JSON(link).Render()
}

func renderHighlightsTable(ctx context.Context, highlights []linkwarden.Highlight) {
	tables.Highlights(ctx, "Highlights", highlights).Render()
}
