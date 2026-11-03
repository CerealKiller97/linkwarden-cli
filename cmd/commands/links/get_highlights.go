package links

import (
	"log/slog"

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
		slog.Error("Failed to resolve link ID", slog.Any("error", err))
		return err
	}

	highlights, err := client.Links.Highlights(cmd.Context(), id)
	if err != nil {
		slog.Error("Failed to get highlights", slog.Any("error", err))
		return err
	}

	renderHighlightsTable(highlights)

	return nil
}
