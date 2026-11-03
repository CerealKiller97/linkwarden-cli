package tags

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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
		slog.Error("Failed to resolve tag ID", slog.Any("error", err))
		return err
	}

	tag, err := client.Tags.Update(cmd.Context(), id, cmd.Flag("name").Value.String())
	if err != nil {
		slog.Error("Failed to update tag", slog.Any("error", err))
		return err
	}

	slog.Info("Tag updated successfully", slog.Int("id", tag.ID), slog.String("name", tag.Name))

	return nil
}
