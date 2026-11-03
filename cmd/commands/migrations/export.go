package migrations

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Export() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export Linkwarden data to a JSON file",
		RunE:  executeExport,
	}

	cmd.Flags().String("output", "", "Path to write the export to. Defaults to ./linkwarden-export-<timestamp>.json.")

	return cmd
}

func executeExport(cmd *cobra.Command, _ []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	if output == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		output = filepath.Join(cwd, "linkwarden-export-"+time.Now().Format("2006-01-02T15-04-05")+".json")
	}

	client := utils.NewLinkwardenClient()

	data, err := client.Migrations.Export(cmd.Context())
	if err != nil {
		slog.Error("Failed to export data", slog.Any("error", err))
		return err
	}

	if err := os.WriteFile(output, data, 0o600); err != nil {
		slog.Error("Failed to write export file", slog.Any("error", err))
		return err
	}

	slog.Info("Exported successfully", slog.String("path", output))

	return nil
}
