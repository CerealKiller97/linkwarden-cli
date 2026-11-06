package highlights

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderHighlightTable(highlight linkwarden.Highlight) {
	tables.Highlights("Highlights", []linkwarden.Highlight{highlight}).Render()
}
