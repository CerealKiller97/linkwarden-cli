package favicon

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "favicon",
		Short: "Fetch favicons",
	}

	cmd.AddCommand(Get(app))

	return cmd
}
