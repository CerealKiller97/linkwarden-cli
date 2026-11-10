package highlights

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create or update a highlight on a link",
		RunE:  executeCreate,
	}

	cmd.Flags().Int("link-id", 0, "Link ID.")
	cmd.Flags().String("text", "", "Highlighted text.")
	cmd.Flags().String("color", "", "Highlight color.")
	cmd.Flags().String("comment", "", "Optional comment.")
	cmd.Flags().Int("start-offset", 0, "Start offset of the highlighted text.")
	cmd.Flags().Int("end-offset", 0, "End offset of the highlighted text.")

	for _, name := range []string{"text", "color", "start-offset", "end-offset"} {
		if err := cmd.MarkFlagRequired(name); err != nil {
			return nil
		}
	}

	return cmd
}

func executeCreate(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	linkID, err := utils.ResolveIDPaged(
		cmd,
		"link-id",
		"Select Link to highlight:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	text, err := cmd.Flags().GetString("text")
	if err != nil {
		return err
	}
	color, err := cmd.Flags().GetString("color")
	if err != nil {
		return err
	}
	comment, err := cmd.Flags().GetString("comment")
	if err != nil {
		return err
	}
	startOffset, err := cmd.Flags().GetInt("start-offset")
	if err != nil {
		return err
	}
	endOffset, err := cmd.Flags().GetInt("end-offset")
	if err != nil {
		return err
	}

	spinner := utils.StartSpinner("Creating highlight…")
	highlight, err := client.Highlights.Create(cmd.Context(), linkwarden.CreateHighlightRequest{
		LinkID:      linkID,
		Text:        text,
		Color:       color,
		Comment:     comment,
		StartOffset: startOffset,
		EndOffset:   endOffset,
	})
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("create highlight: %w", err)
	}

	renderHighlightTable(*highlight)

	return nil
}
