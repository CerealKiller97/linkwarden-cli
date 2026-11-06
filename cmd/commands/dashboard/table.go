package dashboard

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(links []linkwarden.Link) {
	tables.Links("Dashboard", links).Render()
}

func renderDashboardV2Stats(data *linkwarden.DashboardV2Data) {
	t := tables.Links("Dashboard", data.Links)
	t.Caption(fmt.Sprintf("Pinned links: %d · Tags: %d", data.NumberOfPinnedLinks, data.NumberOfTags))
	t.Render()
}
