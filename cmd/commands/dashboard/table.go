package dashboard

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(app *app.App, links []linkwarden.Link) {
	tables.Links(app, "Dashboard", links).Render()
}

func renderDashboardV2Stats(app *app.App, data *linkwarden.DashboardV2Data) {
	t := tables.Links(app, "Dashboard", data.Links)
	t.JSON(data)
	t.Caption(fmt.Sprintf("Pinned links: %d · Tags: %d", data.NumberOfPinnedLinks, data.NumberOfTags))
	t.Render()
}
