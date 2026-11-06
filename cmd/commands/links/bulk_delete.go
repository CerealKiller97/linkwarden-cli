package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func BulkDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-delete",
		Short: "Delete multiple links",
		RunE:  executeBulkDelete,
	}

	cmd.Flags().String("ids", "", "Comma-separated list of link IDs.")

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeBulkDelete(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		return fmt.Errorf("parse --ids: %w", err)
	}

	client := utils.NewLinkwardenClient()

	if err := client.Links.BulkDelete(cmd.Context(), ids); err != nil {
		return fmt.Errorf("delete links: %w", err)
	}

	log.Info().Ints("ids", ids).Msg("Links deleted successfully")

	return nil
}
