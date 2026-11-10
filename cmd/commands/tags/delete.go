package tags

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
		Short: "Delete a tag",
		Example: `  # Pick the tag from a list, then confirm
  linkwarden tags delete

  # Delete tag 14 without asking
  linkwarden tags delete --id 14 --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDelete(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Tag ID.")
	utils.AddYesFlag(cmd)

	return cmd
}

func executeDelete(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		app, cmd,
		"id",
		"Select Tag to delete:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve tag ID: %w", err)
	}

	ok, err := utils.Confirm(app, cmd, fmt.Sprintf("Delete tag %d? This can't be undone.", id))
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was deleted")
		return nil
	}

	spinner := utils.StartSpinner("Deleting tag…")
	err = client.Tags.Delete(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete tag: %w", err)
	}

	log.Info().Int("id", id).Msg("Tag deleted successfully")

	return nil
}
