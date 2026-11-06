package dashboard

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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

	links, err := client.Dashboard.GetV1(cmd.Context())
	if err != nil {
		log.Err(err).Msg("Failed to get dashboard data")
		return err
	}

	renderLinksTable(links)

	return nil
}
