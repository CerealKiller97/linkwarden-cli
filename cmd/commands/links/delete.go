package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a link",
		RunE:  executeDelete,
	}

	cmd.Flags().String("id", "", "Link ID.")

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Link to delete:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	if err := client.Links.Delete(cmd.Context(), id); err != nil {
		return fmt.Errorf("delete link: %w", err)
	}

	log.Info().Int("id", id).Msg("Link deleted successfully")

	return nil
}
