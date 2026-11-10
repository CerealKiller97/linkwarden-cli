package migrations

import (
	"fmt"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/choice"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var importFormats = []choice.Option[linkwarden.ImportFormat]{
	{Name: "linkwarden", Value: linkwarden.ImportFormatLinkwardenJSON, Description: "Linkwarden JSON export"},
	{Name: "html", Value: linkwarden.ImportFormatBookmarksHTML, Description: "Browser bookmarks exported as HTML"},
	{Name: "wallabag", Value: linkwarden.ImportFormatWallabagJSON, Description: "Wallabag JSON export"},
	{Name: "omnivore", Value: linkwarden.ImportFormatOmnivoreJSON, Description: "Omnivore JSON export"},
	{Name: "pocket", Value: linkwarden.ImportFormatPocketJSON, Description: "Pocket JSON export"},
}

func Import() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import Linkwarden data from a data source",
		RunE:  executeImport,
	}

	cmd.Flags().String("file", "", "Path to the file containing the data to import.")
	choice.Flag(cmd, "format", "", "Import source format.", importFormats)

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
	format, _, err := choice.Get[linkwarden.ImportFormat](cmd, "format")
	if err != nil {
		return err
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read import file: %w", err)
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Importing bookmarks…")
	message, err := client.Migrations.Import(cmd.Context(), format, string(fileBytes))
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("import data: %w", err)
	}

	log.Info().Str("message", message).Msg("Import started")

	return nil
}
