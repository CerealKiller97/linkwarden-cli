package tags

import (
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
		"Select Tag to update:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve tag ID")
		return err
	}

	tag, err := client.Tags.Update(cmd.Context(), id, cmd.Flag("name").Value.String())
	if err != nil {
		log.Err(err).Msg("Failed to update tag")
		return err
	}

	log.Info().Int("id", tag.ID).Str("name", tag.Name).Msg("Tag updated successfully")

	return nil
}
