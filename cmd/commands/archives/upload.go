package archives

import (
	"os"
	"path/filepath"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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

	linkID, err := utils.ResolveID(
		cmd.Flag("link-id").Value.String(),
		"--link-id",
		"Select Link:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve link ID")
		return err
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
		log.Err(err).Msg("Failed to open file")
		return err
	}
	defer f.Close()

	link, err := client.Archives.Upload(cmd.Context(), linkID, linkwarden.ArchiveFormat(format), preview, filepath.Base(filePath), f)
	if err != nil {
		log.Err(err).Msg("Failed to upload archive")
		return err
	}

	renderLinkTable(*link)

	return nil
}
