package archives

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinkTable(link linkwarden.Link) {
	t := utils.NewTable("Link")
	t.Columns(
		utils.Column{Header: "ID"},
		utils.Column{Header: "Name", Flex: true},
		utils.Column{Header: "URL", Flex: true},
		utils.Column{Header: "Collection"},
	)
	t.Row(link.ID, link.Name, link.URL, link.Collection.Name)
	t.Render()
}
