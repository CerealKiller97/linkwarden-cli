package config

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get public runtime configuration",
		Example: `  # Show your instance's public configuration
  linkwarden config get

  # Same, as JSON
  linkwarden config get --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGet(app, cmd, args)
		},
	}

	return cmd
}

func executeGet(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching instance config…")
	cfg, err := client.Config.Get(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get configuration: %w", err)
	}

	renderConfigTable(app, cfg)

	return nil
}
