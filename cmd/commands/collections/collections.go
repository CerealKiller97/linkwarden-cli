package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collections",
		Short: "Manage collections",
	}

	cmd.AddCommand(Create(app))
	cmd.AddCommand(List(app))
	cmd.AddCommand(Get(app))
	cmd.AddCommand(Update(app))
	cmd.AddCommand(Delete(app))

	return cmd
}
