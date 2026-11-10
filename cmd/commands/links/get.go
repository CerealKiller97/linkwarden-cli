package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a link by ID",
		Example: `  # Pick the link from a list
  linkwarden links get

  # Get link 412
  linkwarden links get --id 412

  # Print only its URL
  linkwarden links get --id 412 --json | jq -r .url`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGet(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Link ID.")

	return cmd
}

func executeGet(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveIDPaged(
		app, cmd,
		"id",
		"Select Link:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	spinner := utils.StartSpinner("Fetching link…")
	link, err := client.Links.Get(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get link: %w", err)
	}

	renderLinkTable(app, link)

	return nil
}
