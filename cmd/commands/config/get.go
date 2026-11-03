package config

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get public runtime configuration",
		RunE:  executeGet,
	}

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	cfg, err := client.Config.Get(cmd.Context())
	if err != nil {
		slog.Error("Failed to get configuration", slog.Any("error", err))
		return err
	}

	renderConfigTable(cfg)

	return nil
}
