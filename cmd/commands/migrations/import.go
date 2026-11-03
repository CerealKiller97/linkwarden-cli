package migrations

import (
	"log/slog"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

const importFormatHelp = "Import source format: 0=Linkwarden JSON, 1=Bookmarks HTML, " +
	"2=Wallabag JSON, 3=Omnivore JSON, 4=Pocket JSON."

func Import() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import Linkwarden data from a data source",
		RunE:  executeImport,
	}

	cmd.Flags().String("file", "", "Path to the file containing the data to import.")
	cmd.Flags().Int("format", 0, importFormatHelp)

	for _, name := range []string{"file", "format"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return nil
		}
	}

	return cmd
}

func executeImport(cmd *cobra.Command, _ []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	format, err := cmd.Flags().GetInt("format")
	if err != nil {
		return err
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Failed to read import file", slog.Any("error", err))
		return err
	}

	client := utils.NewLinkwardenClient()

	message, err := client.Migrations.Import(cmd.Context(), linkwarden.ImportFormat(format), string(fileBytes))
	if err != nil {
		slog.Error("Failed to import data", slog.Any("error", err))
		return err
	}

	slog.Info("Import started", slog.String("message", message))

	return nil
}
