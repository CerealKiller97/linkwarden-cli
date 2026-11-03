package links

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func BulkDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-delete",
		Short: "Delete multiple links",
		RunE:  executeBulkDelete,
	}

	cmd.Flags().String("ids", "", "Comma-separated list of link IDs.")

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeBulkDelete(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		slog.Error("Failed to parse --ids", slog.Any("error", err))
		return err
	}

	client := utils.NewLinkwardenClient()

	if err := client.Links.BulkDelete(cmd.Context(), ids); err != nil {
		slog.Error("Failed to delete links", slog.Any("error", err))
		return err
	}

	slog.Info("Links deleted successfully", slog.Any("ids", ids))

	return nil
}
