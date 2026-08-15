package links

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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
		"Select Link:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve link ID")
		return err
	}

	highlights, err := client.Links.Highlights(cmd.Context(), id)
	if err != nil {
		log.Err(err).Msg("Failed to get highlights")
		return err
	}

	renderHighlightsTable(highlights)

	return nil
}
