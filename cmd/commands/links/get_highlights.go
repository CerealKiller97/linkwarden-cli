package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func GetHighlights() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-highlights",
		Short: "Get highlights for a link",
		RunE:  executeGetHighlights,
	}

	cmd.Flags().String("id", "", "Link ID.")

	return cmd
}

func executeGetHighlights(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveIDPaged(
		cmd.Flag("id").Value.String(),
		"--id",
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

	renderHighlightsTable(highlights)

	return nil
}
