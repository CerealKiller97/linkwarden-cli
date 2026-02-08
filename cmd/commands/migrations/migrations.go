package migrations

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/services/migration"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var (
	importFormat string
	exportPath   string
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrations",
		Short: "Manage migration (Import / Export)",
	}

	cmd.AddCommand(Import())
	cmd.AddCommand(Export())
	return cmd
}

func Import() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import Linkwarden data from a data source",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringVarP(
		&importFormat,
		"format",
		"f",
		"",
		"Import format. Valid options are: linkwarden, htmlFile, wallbag or omnivore",
	)

	return cmd
}

func Export() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export Linkwarden data to a JSON file.",

		RunE: func(cmd *cobra.Command, args []string) error {
			collectionService := migration.New(
				utils.GetInstanceURL(),
				utils.GetLinkwardenToken(),
				zerolog.New(os.Stdout),
			)

			data, backupFileName, err := collectionService.Export(cmd.Context())
			if err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			filePath := path.Join(cwd, backupFileName)
			f, err := os.Create(filePath)

			_, err = f.Write(data)
			if err != nil {
				return err
			}

			log.Info().Str("path", filePath).Msg("Exported successfully")

			return nil
		},
	}

	//cmd.Flags().StringVarP(&exportFilePath, "source", "s", "", "Source directory to read from")

	return cmd
}
