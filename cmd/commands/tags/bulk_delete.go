package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func BulkDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-delete",
		Short: "Delete multiple tags",
		RunE:  executeBulkDelete,
	}

	cmd.Flags().String("ids", "", "Comma-separated list of tag IDs.")

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeBulkDelete(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		log.Err(err).Msg("Failed to parse --ids")
		return err
	}

	client := utils.NewLinkwardenClient()

	if err := client.Tags.BulkDelete(cmd.Context(), ids); err != nil {
		log.Err(err).Msg("Failed to delete tags")
		return err
	}

	log.Info().Ints("ids", ids).Msg("Tags deleted successfully")

	return nil
}
