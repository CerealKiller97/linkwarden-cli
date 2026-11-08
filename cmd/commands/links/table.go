package links

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinksTable(links []linkwarden.Link) {
	tables.Links("Links", links).Render()
}

// renderLinkTable prints one link, as a JSON object with --json.
func renderLinkTable(link *linkwarden.Link) {
	tables.Links("Links", []linkwarden.Link{*link}).JSON(link).Render()
}

func renderHighlightsTable(highlights []linkwarden.Highlight) {
	tables.Highlights("Highlights", highlights).Render()
}
