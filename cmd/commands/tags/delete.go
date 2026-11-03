package tags

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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
		"--id",
		"Select Tag to delete:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		slog.Error("Failed to resolve tag ID", slog.Any("error", err))
		return err
	}

	if err := client.Tags.Delete(cmd.Context(), id); err != nil {
		slog.Error("Failed to delete tag", slog.Any("error", err))
		return err
	}

	slog.Info("Tag deleted successfully", slog.Int("id", id))

	return nil
}
