package links

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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
		slog.Error("Failed to resolve link ID", slog.Any("error", err))
		return err
	}

	if err := client.Links.Delete(cmd.Context(), id); err != nil {
		slog.Error("Failed to delete link", slog.Any("error", err))
		return err
	}

	slog.Info("Link deleted successfully", slog.Int("id", id))

	return nil
}
