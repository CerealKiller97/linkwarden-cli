package rss

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an RSS subscription",
		RunE:  executeDelete,
	}

	cmd.Flags().String("id", "", "Subscription ID.")

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select RSS subscription to delete:",
		rssPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		slog.Error("Failed to resolve RSS subscription ID", slog.Any("error", err))
		return err
	}

	if err := client.RSS.Delete(cmd.Context(), id); err != nil {
		slog.Error("Failed to delete RSS subscription", slog.Any("error", err))
		return err
	}

	slog.Info("RSS subscription deleted successfully", slog.Int("id", id))

	return nil
}
