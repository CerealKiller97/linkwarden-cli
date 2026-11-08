package archives

import (
	"fmt"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const formatHelp = "Archive format: 0=PNG, 1=JPEG, 2=PDF, 3=Readability JSON, 4=Monolith HTML."

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve a link's archive file",
		RunE:  executeGet,
	}

	cmd.Flags().String("link-id", "", "Link ID.")
	cmd.Flags().Int("format", 0, formatHelp)
	cmd.Flags().Bool("preview", false, "Retrieve the preview instead of the main archive file.")
	cmd.Flags().String("output", "", "Path to save the archive file to. Without it, only content-type/size are printed.")

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	linkID, err := utils.ResolveID(
		cmd.Flag("link-id").Value.String(),
		"--link-id",
		"Select Link:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	var format linkwarden.ArchiveFormat
	if cmd.Flags().Changed("format") {
		f, err := cmd.Flags().GetInt("format")
		if err != nil {
			return err
		}
		format = linkwarden.ArchiveFormat(f)
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
		fmt.Printf("Content-Type: %s, Size: %d bytes\n", file.ContentType, len(file.Data))
		return nil
	}

	if err := os.WriteFile(output, file.Data, 0o600); err != nil {
		return fmt.Errorf("write archive file: %w", err)
	}

	log.Info().Str("path", output).Msg("Archive saved successfully")

	return nil
}
