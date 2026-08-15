package config

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Msg("Failed to get configuration")
		return err
	}

	renderConfigTable(cfg)

	return nil
}
