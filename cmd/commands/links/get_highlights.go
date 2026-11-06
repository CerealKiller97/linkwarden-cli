package links

import (
	"fmt"
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

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Link:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	highlights, err := client.Links.Highlights(cmd.Context(), id)
	if err != nil {
		return fmt.Errorf("get highlights: %w", err)
	}

	renderHighlightsTable(highlights)

	return nil
}
