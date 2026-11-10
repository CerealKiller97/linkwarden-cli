package highlights

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderHighlightTable(ctx context.Context, highlight linkwarden.Highlight) {
	tables.Highlights(ctx, "Highlights", []linkwarden.Highlight{highlight}).JSON(highlight).Render()
}
