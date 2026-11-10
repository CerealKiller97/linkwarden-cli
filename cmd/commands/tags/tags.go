package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage tags",
	}

	cmd.AddCommand(List(app))
	cmd.AddCommand(Get(app))
	cmd.AddCommand(Create(app))
	cmd.AddCommand(Update(app))
	cmd.AddCommand(Delete(app))
	cmd.AddCommand(BulkDelete(app))

	return cmd
}
