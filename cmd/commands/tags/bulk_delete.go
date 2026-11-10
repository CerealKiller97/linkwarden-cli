package tags

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func BulkDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-delete",
		Short: "Delete multiple tags",
		Example: `  # Delete three tags (asks first)
  linkwarden tags bulk-delete --ids 14,15,22

  # Same, without asking
  linkwarden tags bulk-delete --ids 14 --ids 15 --ids 22 --yes`,
		RunE: executeBulkDelete,
	}

	utils.AddIDsFlag(cmd, "Tag IDs.")
	utils.AddYesFlag(cmd)

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeBulkDelete(cmd *cobra.Command, _ []string) error {
	ids, err := utils.GetIDs(cmd)
	if err != nil {
		return err
	}

	ok, err := utils.Confirm(cmd, fmt.Sprintf("Delete %s (%s)? This can't be undone.",
		utils.Plural(len(ids), "tag"), utils.JoinIDs(ids)))
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was deleted")
		return nil
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Deleting tags…")
	err = client.Tags.BulkDelete(cmd.Context(), ids)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete tags: %w", err)
	}

	log.Info().Ints("ids", ids).Msg("Tags deleted successfully")

	return nil
}
