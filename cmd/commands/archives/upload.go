package archives

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Upload() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a local file as a link's archive",
		RunE:  executeUpload,
	}

	cmd.Flags().String("link-id", "", "Link ID.")
	cmd.Flags().String("file", "", "Local file to upload.")
	cmd.Flags().Int("format", 0, formatHelp)
	cmd.Flags().Bool("preview", false, "Store the uploaded file as the generated preview instead of the main archive.")

	for _, name := range []string{"file", "format"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return nil
		}
	}

	return cmd
}

func executeUpload(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	linkID, err := utils.ResolveIDPaged(
		cmd.Flag("link-id").Value.String(),
		"--link-id",
		"Select Link:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	filePath, err := cmd.Flags().GetString("file")
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
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	spinner := utils.StartSpinner("Uploading archive…")
	link, err := client.Archives.Upload(cmd.Context(), linkID, linkwarden.ArchiveFormat(format), preview, filepath.Base(filePath), f)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("upload archive: %w", err)
	}

	renderLinkTable(*link)

	return nil
}
