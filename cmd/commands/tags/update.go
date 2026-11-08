package tags

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Update() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Rename a tag",
		RunE:  executeUpdate,
	}

	cmd.Flags().String("id", "", "Tag ID.")
	cmd.Flags().String("name", "", "New tag name.")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}

	return cmd
}

func executeUpdate(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Tag to update:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve tag ID: %w", err)
	}

	spinner := utils.StartSpinner("Updating tag…")
	tag, err := client.Tags.Update(cmd.Context(), id, cmd.Flag("name").Value.String())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("update tag: %w", err)
	}

	log.Info().Int("id", tag.ID).Str("name", tag.Name).Msg("Tag updated successfully")

	return nil
}
