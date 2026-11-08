package migrations

import (
	"fmt"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
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
		return fmt.Errorf("read import file: %w", err)
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Importing bookmarks…")
	message, err := client.Migrations.Import(cmd.Context(), linkwarden.ImportFormat(format), string(fileBytes))
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("import data: %w", err)
	}

	log.Info().Str("message", message).Msg("Import started")

	return nil
}
