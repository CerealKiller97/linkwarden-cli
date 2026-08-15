package dashboard

import (
	"fmt"
	"os"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/jedib0t/go-pretty/v6/table"
)

func renderLinksTable(title string, links []linkwarden.Link) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(title)
	t.AppendHeader(table.Row{"ID", "Name", "URL", "Collection", "Tags"})

	for _, l := range links {
		tagNames := make([]string, len(l.Tags))
		for i, tag := range l.Tags {
			tagNames[i] = tag.Name
		}

		t.AppendRow(table.Row{l.ID, l.Name, l.URL, l.Collection.Name, strings.Join(tagNames, ", ")})
	}

	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
	t.Render()
}

func renderDashboardV2Stats(data *linkwarden.DashboardV2Data) {
	fmt.Printf("Pinned links: %d, Tags: %d\n", data.NumberOfPinnedLinks, data.NumberOfTags)
	renderLinksTable("Dashboard (v2)", data.Links)
}
