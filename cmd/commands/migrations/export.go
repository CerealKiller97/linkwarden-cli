package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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

	spinner := utils.StartSpinner("Exporting data…")
	data, err := client.Migrations.Export(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("export data: %w", err)
	}

	if err := os.WriteFile(output, data, 0o600); err != nil {
		return fmt.Errorf("write export file: %w", err)
	}

	log.Info().Str("path", output).Msg("Exported successfully")

	return nil
}
