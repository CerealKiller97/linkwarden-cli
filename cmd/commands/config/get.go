package config

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get public runtime configuration",
		Example: `  # Show your instance's public configuration
  linkwarden config get

  # Same, as JSON
  linkwarden config get --json`,
		RunE: executeGet,
	}

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching instance config…")
	cfg, err := client.Config.Get(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get configuration: %w", err)
	}

	renderConfigTable(cfg)

	return nil
}
