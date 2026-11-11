package tokens

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Manage tokens",
	}

	cmd.AddCommand(List(app))
	cmd.AddCommand(Generate(app))
	cmd.AddCommand(Revoke(app))

	return cmd
}
