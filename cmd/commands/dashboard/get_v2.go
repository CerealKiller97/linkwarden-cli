package dashboard

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func GetV2() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-v2",
		Short: "Get dashboard data (v2)",
		RunE:  executeGetV2,
	}

	return cmd
}

func executeGetV2(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching dashboard…")
	data, err := client.Dashboard.GetV2(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get dashboard data: %w", err)
	}

	renderDashboardV2Stats(data)

	return nil
}
