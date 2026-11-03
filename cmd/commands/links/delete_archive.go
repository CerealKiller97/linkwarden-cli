package links

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func DeleteArchive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-archive",
		Short: "Delete archived files for multiple links",
		RunE:  executeDeleteArchive,
	}

	cmd.Flags().String("ids", "", "Comma-separated list of link IDs.")

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeDeleteArchive(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		slog.Error("Failed to parse --ids", slog.Any("error", err))
		return err
	}

	client := utils.NewLinkwardenClient()

	if err := client.Links.DeleteArchive(cmd.Context(), ids); err != nil {
		slog.Error("Failed to delete archived files", slog.Any("error", err))
		return err
	}

	slog.Info("Archived files deleted successfully", slog.Any("ids", ids))

	return nil
}
