package highlights

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a highlight",
		Example: `  # Pick a link, then one of its highlights
  linkwarden highlights delete

  # Pick from the highlights on link 412
  linkwarden highlights delete --link-id 412

  # Delete highlight 31 without asking
  linkwarden highlights delete --id 31 --yes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDelete(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Highlight ID.")
	cmd.Flags().Int("link-id", 0, "Link ID (used to pick a highlight when --id is not set).")
	utils.AddYesFlag(cmd)

	return cmd
}

func executeDelete(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := resolveHighlightID(app, cmd, client)
	if err != nil {
		return fmt.Errorf("resolve highlight ID: %w", err)
	}

	ok, err := utils.Confirm(app, cmd, fmt.Sprintf("Delete highlight %d? This can't be undone.", id))
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was deleted")
		return nil
	}

	spinner := utils.StartSpinner("Deleting highlight…")
	err = client.Highlights.Delete(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete highlight: %w", err)
	}

	log.Info().Int("id", id).Msg("Highlight deleted successfully")

	return nil
}

// resolveHighlightID returns --id when set; otherwise it walks the user through a two-step picker (link, then one of its highlights), since
// there is no list-all-highlights endpoint to pick from directly.
func resolveHighlightID(app *app.App, cmd *cobra.Command, client *linkwarden.Client) (int, error) {
	if cmd.Flags().Changed("id") {
		return utils.ResolveID(app, cmd, "id", "", nil)
	}

	linkID, err := utils.ResolveIDPaged(
		app, cmd,
		"link-id",
		"Select Link:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return 0, err
	}

	// --id is unset here, so this opens the highlight picker.
	return utils.ResolveID(app, cmd, "id", "Select Highlight to delete:", highlightsPickerFetch(cmd.Context(), client, linkID))
}
