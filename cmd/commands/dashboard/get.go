package dashboard

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get dashboard data (v1)",
		RunE:  executeGet,
	}

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching dashboard…")
	links, err := client.Dashboard.GetV1(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get dashboard data: %w", err)
	}

	renderLinksTable(links)

	return nil
}
