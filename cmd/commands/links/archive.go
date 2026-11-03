package links

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Archive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archive",
		Short: "Trigger archiving of a link",
		RunE:  executeArchive,
	}

	cmd.Flags().String("id", "", "Link ID.")

	return cmd
}

func executeArchive(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Link to archive:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		slog.Error("Failed to resolve link ID", slog.Any("error", err))
		return err
	}

	if err := client.Links.Archive(cmd.Context(), id); err != nil {
		slog.Error("Failed to archive link", slog.Any("error", err))
		return err
	}

	slog.Info("Link archiving triggered", slog.Int("id", id))

	return nil
}
