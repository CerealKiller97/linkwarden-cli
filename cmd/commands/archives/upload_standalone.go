package archives

import (
	"os"
	"path/filepath"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func UploadStandalone() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-standalone",
		Short: "Upload a local file as a new link's archive",
		RunE:  executeUploadStandalone,
	}

	cmd.Flags().String("file", "", "Local file to upload.")
	cmd.Flags().String("url", "", "URL to associate with the new link.")
	cmd.Flags().Int("format", 0, formatHelp)
	cmd.Flags().Bool("preview", false, "Store the uploaded file as the generated preview instead of the main archive.")

	for _, name := range []string{"file", "url", "format"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return nil
		}
	}

	return cmd
}

func executeUploadStandalone(cmd *cobra.Command, _ []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	linkURL, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}
	format, err := cmd.Flags().GetInt("format")
	if err != nil {
		return err
	}
	preview, err := cmd.Flags().GetBool("preview")
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Err(err).Msg("Failed to open file")
		return err
	}
	defer f.Close()

	client := utils.NewLinkwardenClient()

	link, err := client.Archives.UploadStandalone(cmd.Context(), linkwarden.ArchiveFormat(format), preview, linkURL, filepath.Base(filePath), f)
	if err != nil {
		log.Err(err).Msg("Failed to upload standalone archive")
		return err
	}

	renderLinkTable(*link)

	return nil
}
