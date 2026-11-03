package collections

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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
		slog.Error("Failed to resolve collection ID", slog.Any("error", err))
		return err
	}

	if err := client.Collections.Delete(cmd.Context(), id); err != nil {
		slog.Error("Failed to delete collection", slog.Any("error", err))
		return err
	}

	slog.Info("Collection deleted successfully", slog.Int("id", id))

	return nil
}
