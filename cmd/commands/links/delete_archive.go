package links

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func DeleteArchive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-archive",
		Short: "Delete archived files for multiple links",
		RunE:  executeDeleteArchive,
	}

	cmd.Flags().String("ids", "", "Comma-separated list of link IDs.")

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeDeleteArchive(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		log.Err(err).Msg("Failed to parse --ids")
		return err
	}

	client := utils.NewLinkwardenClient()

	if err := client.Links.DeleteArchive(cmd.Context(), ids); err != nil {
		log.Err(err).Msg("Failed to delete archived files")
		return err
	}

	log.Info().Ints("ids", ids).Msg("Archived files deleted successfully")

	return nil
}
