package dashboard

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(title string, links []linkwarden.Link) {
	utils.RenderLinksTable(title, links)
}

func renderDashboardV2Stats(data *linkwarden.DashboardV2Data) {
	fmt.Printf("Pinned links: %d, Tags: %d\n", data.NumberOfPinnedLinks, data.NumberOfTags)
	renderLinksTable("Dashboard (v2)", data.Links)
}
