package archives

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
)

func renderLinkTable(link linkwarden.Link) {
	t := utils.NewTable("Link")
	t.Header("ID", "Name", "URL", "Collection")
	t.Row(link.ID, link.Name, link.URL, link.Collection.Name)
	t.Render()
}
