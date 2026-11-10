package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func GetHighlights(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-highlights",
		Short: "Get highlights for a link",
		Example: `  # Show the highlights on link 412
  linkwarden links get-highlights --id 412

  # Pick the link from a list, and print JSON
  linkwarden links get-highlights --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGetHighlights(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Link ID.")

	return cmd
}

func executeGetHighlights(app *app.App, cmd *cobra.Command, _ []string) error {
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

	spinner := utils.StartSpinner("Fetching highlights…")
	highlights, err := client.Links.Highlights(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get highlights: %w", err)
	}

	renderHighlightsTable(app, highlights)

	return nil
}
