package archives

import (
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/tables"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinkTable(link linkwarden.Link) {
	tables.Links("Links", []linkwarden.Link{link}).JSON(link).Render()
}
