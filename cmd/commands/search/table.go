package search

import (
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderSearchResultsTable(links []linkwarden.Link) {
	t := utils.NewTable("Search Results")
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
