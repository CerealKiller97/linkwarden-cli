package archives

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/choice"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
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
	choice.Flag(cmd, "format", "", "Archive format.", uploadFormats)
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
	format, _, err := choice.Get[linkwarden.ArchiveFormat](cmd, "format")
	if err != nil {
		return err
	}
	preview, err := cmd.Flags().GetBool("preview")
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Uploading archive…")
	link, err := client.Archives.UploadStandalone(cmd.Context(), format, preview, linkURL, filepath.Base(filePath), f)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("upload standalone archive: %w", err)
	}

	renderLinkTable(*link)

	return nil
}
