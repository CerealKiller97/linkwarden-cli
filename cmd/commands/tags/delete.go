package tags

import (
	"fmt"

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
	utils.AddYesFlag(cmd)

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Tag to delete:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve tag ID: %w", err)
	}

	ok, err := utils.Confirm(cmd, fmt.Sprintf("Delete tag %d? This can't be undone.", id))
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
