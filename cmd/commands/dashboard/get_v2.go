package dashboard

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func GetV2(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-v2",
		Short: "Get dashboard data (v2)",
		Example: `  # Show the v2 dashboard data
  linkwarden dashboard get-v2

  # Same, as JSON
  linkwarden dashboard get-v2 --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGetV2(app, cmd, args)
		},
	}

	return cmd
}

func executeGetV2(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching dashboard…")
	data, err := client.Dashboard.GetV2(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get dashboard data: %w", err)
	}

	renderDashboardV2Stats(app, data)

	return nil
}
