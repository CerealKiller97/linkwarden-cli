package rss

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func List(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List RSS subscriptions",
		Example: `  # List your feed subscriptions
  linkwarden rss list

  # Same, as JSON
  linkwarden rss list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeList(app, cmd, args)
		},
	}

	return cmd
}

func executeList(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching RSS subscriptions…")
	subs, err := client.RSS.List(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("list RSS subscriptions: %w", err)
	}

	renderRSSTable(app, subs, subs)

	return nil
}
