package search

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderSearchResultsTable(links []linkwarden.Link) {
	utils.RenderLinksTable("Search Results", links)
}
