package highlights

import (
	"fmt"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"strconv"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a highlight",
		RunE:  executeDelete,
	}

	cmd.Flags().String("id", "", "Highlight ID.")
	cmd.Flags().String("link-id", "", "Link ID (used to pick a highlight when --id is not set).")

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := resolveHighlightID(cmd, client, cmd.Flag("id").Value.String())
	if err != nil {
		return fmt.Errorf("resolve highlight ID: %w", err)
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

// resolveHighlightID returns idFlag parsed as an int when non-empty; otherwise it
// walks the user through a two-step picker (link, then one of its highlights), since
// there is no list-all-highlights endpoint to pick from directly.
func resolveHighlightID(cmd *cobra.Command, client *linkwarden.Client, idFlag string) (int, error) {
	if idFlag != "" {
		return strconv.Atoi(idFlag)
	}

	linkID, err := utils.ResolveIDPaged(
		cmd.Flag("link-id").Value.String(),
		"--link-id",
		"Select Link:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return 0, err
	}

	return utils.ResolveID("", "--id", "Select Highlight to delete:", highlightsPickerFetch(cmd.Context(), client, linkID))
}
