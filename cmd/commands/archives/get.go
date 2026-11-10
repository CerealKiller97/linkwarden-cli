package archives

import (
	"fmt"
	"os"
	"slices"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/choice"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// uploadFormats are the archive formats a file can be uploaded as.
var uploadFormats = []choice.Option[linkwarden.ArchiveFormat]{
	{Name: "png", Value: linkwarden.ArchiveFormatPNG, Description: "PNG screenshot"},
	{Name: "jpeg", Value: linkwarden.ArchiveFormatJPEG, Description: "JPEG screenshot"},
	{Name: "pdf", Value: linkwarden.ArchiveFormatPDF, Description: "PDF"},
	{Name: "html", Value: linkwarden.ArchiveFormatMonolithHTML, Description: "Single-file HTML page (Monolith)"},
}

// getFormats are the archive formats that can be retrieved: the upload formats
// plus the Readability text Linkwarden extracts itself.
var getFormats = append(slices.Clone(uploadFormats),
	choice.Option[linkwarden.ArchiveFormat]{Name: "readable", Value: linkwarden.ArchiveFormatReadabilityJSON, Description: "Readability article text (JSON)"},
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve a link's archive file",
		RunE:  executeGet,
	}

	cmd.Flags().String("link-id", "", "Link ID.")
	choice.Flag(cmd, "format", "png", "Archive format.", getFormats)
	cmd.Flags().Bool("preview", false, "Retrieve the preview instead of the main archive file.")
	cmd.Flags().String("output", "", "Path to save the archive file to. Without it, only content-type/size are printed.")

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
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

	format, _, err := choice.Get[linkwarden.ArchiveFormat](cmd, "format")
	if err != nil {
		return err
	}

	preview, err := cmd.Flags().GetBool("preview")
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	spinner := utils.StartSpinner("Downloading archive…")
	file, err := client.Archives.Get(cmd.Context(), linkID, format, preview)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get archive: %w", err)
	}

	if output == "" {
		return utils.PrintFileInfo(file.ContentType, len(file.Data))
	}

	if err := os.WriteFile(output, file.Data, 0o600); err != nil {
		return fmt.Errorf("write archive file: %w", err)
	}

	log.Info().Str("path", output).Msg("Archive saved successfully")

	return nil
}
