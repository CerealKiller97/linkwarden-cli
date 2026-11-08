package collections

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a collection",
		RunE:  executeDelete,
	}

	cmd.Flags().String("id", "", "Collection ID.")

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Collection to delete:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve collection ID: %w", err)
	}

	spinner := utils.StartSpinner("Deleting collection…")
	err = client.Collections.Delete(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	log.Info().Int("id", id).Msg("Collection deleted successfully")

	return nil
}
