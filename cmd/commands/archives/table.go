package archives

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinkTable(ctx context.Context, link linkwarden.Link) {
	tables.Links(ctx, "Links", []linkwarden.Link{link}).JSON(link).Render()
}
