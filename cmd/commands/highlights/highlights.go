package highlights

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "highlights",
		Short: "Manage link highlights",
	}

	cmd.AddCommand(Create(app))
	cmd.AddCommand(Delete(app))

	return cmd
}
