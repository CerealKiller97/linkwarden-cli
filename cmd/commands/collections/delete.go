package collections

import (
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
		"Select Collection to delete:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve collection ID")
		return err
	}

	if err := client.Collections.Delete(cmd.Context(), id); err != nil {
		log.Err(err).Msg("Failed to delete collection")
		return err
	}

	log.Info().Int("id", id).Msg("Collection deleted successfully")

	return nil
}
