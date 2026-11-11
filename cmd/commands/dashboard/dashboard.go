package dashboard

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "View and configure the dashboard",
	}

	cmd.AddCommand(Get(app))
	cmd.AddCommand(GetV2(app))
	cmd.AddCommand(UpdateV2(app))

	return cmd
}
