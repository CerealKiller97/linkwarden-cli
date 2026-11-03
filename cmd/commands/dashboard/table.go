package dashboard

import (
	"fmt"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(title string, links []linkwarden.Link) {
	t := utils.NewTable(title)
	t.Header("ID", "Name", "URL", "Collection", "Tags")

	for _, l := range links {
		tagNames := make([]string, len(l.Tags))
		for i, tag := range l.Tags {
			tagNames[i] = tag.Name
		}

		t.Row(l.ID, l.Name, l.URL, l.Collection.Name, strings.Join(tagNames, ", "))
	}

	t.Render()
}

func renderDashboardV2Stats(data *linkwarden.DashboardV2Data) {
	fmt.Printf("Pinned links: %d, Tags: %d\n", data.NumberOfPinnedLinks, data.NumberOfTags)
	renderLinksTable("Dashboard (v2)", data.Links)
}
