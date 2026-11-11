package dashboard

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get dashboard data (v1)",
		Example: `  # Show the dashboard data
  linkwarden dashboard get

  # Same, as JSON
  linkwarden dashboard get --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGet(app, cmd, args)
		},
	}

	return cmd
}

func executeGet(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching dashboard…")
	links, err := client.Dashboard.GetV1(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get dashboard data: %w", err)
	}

	renderLinksTable(app, links)

	return nil
}
