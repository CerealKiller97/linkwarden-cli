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
	utils.AddYesFlag(cmd)

	cobra.CheckErr(cmd.MarkFlagRequired("ids"))

	return cmd
}

func executeBulkDelete(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		return fmt.Errorf("parse --ids: %w", err)
	}

	ok, err := utils.Confirm(cmd, fmt.Sprintf("Delete %s (%s)? This can't be undone.",
		utils.Plural(len(ids), "link"), utils.JoinIDs(ids)))
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was deleted")
		return nil
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Deleting links…")
	err = client.Links.BulkDelete(cmd.Context(), ids)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete links: %w", err)
	}

	log.Info().Ints("ids", ids).Msg("Links deleted successfully")

	return nil
}
