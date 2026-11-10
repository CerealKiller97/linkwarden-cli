package dashboard

import (
	"context"
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(ctx context.Context, links []linkwarden.Link) {
	tables.Links(ctx, "Dashboard", links).Render()
}

func renderDashboardV2Stats(ctx context.Context, data *linkwarden.DashboardV2Data) {
	t := tables.Links(ctx, "Dashboard", data.Links)
	t.JSON(data)
	t.Caption(fmt.Sprintf("Pinned links: %d · Tags: %d", data.NumberOfPinnedLinks, data.NumberOfTags))
	t.Render()
}
