package links

import (
	"fmt"

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
		return fmt.Errorf("parse --ids: %w", err)
	}

	client := utils.NewLinkwardenClient()

	if err := client.Links.DeleteArchive(cmd.Context(), ids); err != nil {
		return fmt.Errorf("delete archived files: %w", err)
	}

	log.Info().Ints("ids", ids).Msg("Archived files deleted successfully")

	return nil
}
