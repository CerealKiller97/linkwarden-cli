package links

import (
	"fmt"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Archive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Trigger archiving of a link",
		RunE:  executeArchive,
	}

	cmd.Flags().String("id", "", "Link ID.")

	return cmd
}

func executeArchive(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveIDPaged(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Link to archive:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	spinner := utils.StartSpinner("Archiving link…")
	err = client.Links.Archive(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("archive link: %w", err)
	}

	log.Info().Int("id", id).Msg("Link archiving triggered")

	return nil
}
