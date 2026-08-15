package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a tag",
		RunE:  executeDelete,
	}

	cmd.Flags().String("id", "", "Tag ID.")

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"Select Tag to delete:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve tag ID")
		return err
	}

	if err := client.Tags.Delete(cmd.Context(), id); err != nil {
		log.Err(err).Msg("Failed to delete tag")
		return err
	}

	log.Info().Int("id", id).Msg("Tag deleted successfully")

	return nil
}
