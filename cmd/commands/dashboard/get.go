package dashboard

import (
	"log/slog"

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

	links, err := client.Dashboard.GetV1(cmd.Context())
	if err != nil {
		slog.Error("Failed to get dashboard data", slog.Any("error", err))
		return err
	}

	renderLinksTable("Dashboard (v1)", links)

	return nil
}
