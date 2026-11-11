package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func BulkDelete(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-delete",
		Short: "Delete multiple links",
		Example: `  # Delete two links (asks first)
  linkwarden links bulk-delete --ids 301,302

  # Same, without asking (scripts, cron)
  linkwarden links bulk-delete --ids 301 --ids 302 --yes

  # Delete every link that matches a search
  linkwarden links list --search "old draft" -q | paste -sd, - | xargs -I{} linkwarden links bulk-delete --ids {} --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeBulkDelete(app, cmd, args)
		},
	}

	utils.AddIDsFlag(cmd, "Link IDs.")
	utils.AddYesFlag(cmd)

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeBulkDelete(app *app.App, cmd *cobra.Command, _ []string) error {
	ids, err := utils.GetIDs(cmd)
	if err != nil {
		return err
	}

	ok, err := utils.Confirm(app, cmd, fmt.Sprintf("Delete %s (%s)? This can't be undone.",
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
