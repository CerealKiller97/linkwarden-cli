package links

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "Manage links",
	}

	cmd.AddCommand(Create(app))
	cmd.AddCommand(Import(app))
	cmd.AddCommand(List(app))
	cmd.AddCommand(Get(app))
	cmd.AddCommand(Update(app))
	cmd.AddCommand(Delete(app))
	cmd.AddCommand(Archive(app))
	cmd.AddCommand(BulkUpdate(app))
	cmd.AddCommand(BulkDelete(app))
	cmd.AddCommand(DeleteArchive(app))
	cmd.AddCommand(GetHighlights(app))

	return cmd
}
