package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a link",
		Example: `  # Pick the link from a list, then confirm
  linkwarden links delete

  # Delete link 412 without asking (scripts, cron)
  linkwarden links delete --id 412 --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDelete(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Link ID.")
	utils.AddYesFlag(cmd)

	return cmd
}

func executeDelete(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveIDPaged(
		app, cmd,
		"id",
		"Select Link to delete:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	ok, err := utils.Confirm(app, cmd, fmt.Sprintf("Delete link %d? This can't be undone.", id))
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was deleted")
		return nil
	}

	spinner := utils.StartSpinner("Deleting link…")
	err = client.Links.Delete(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	log.Info().Int("id", id).Msg("Link deleted successfully")

	return nil
}
