package dashboard

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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

	data, err := client.Dashboard.GetV2(cmd.Context())
	if err != nil {
		log.Err(err).Msg("Failed to get dashboard data")
		return err
	}

	renderDashboardV2Stats(data)

	return nil
}
