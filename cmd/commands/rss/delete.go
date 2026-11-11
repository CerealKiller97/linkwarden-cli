package rss

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an RSS subscription",
		Example: `  # Pick the subscription from a list, then confirm
  linkwarden rss delete

  # Delete subscription 5 without asking
  linkwarden rss delete --id 5 --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDelete(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Subscription ID.")
	utils.AddYesFlag(cmd)

	return cmd
}

func executeDelete(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		app, cmd,
		"id",
		"Select RSS subscription to delete:",
		rssPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve RSS subscription ID: %w", err)
	}

	ok, err := utils.Confirm(app, cmd, fmt.Sprintf("Delete RSS subscription %d? This can't be undone.", id))
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was deleted")
		return nil
	}

	spinner := utils.StartSpinner("Deleting RSS subscription…")
	err = client.RSS.Delete(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete RSS subscription: %w", err)
	}

	log.Info().Int("id", id).Msg("RSS subscription deleted successfully")

	return nil
}
