package migrations

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrations",
		Short: "Manage migration (Import / Export)",
	}

	cmd.AddCommand(Import(app))
	cmd.AddCommand(Export(app))

	return cmd
}
