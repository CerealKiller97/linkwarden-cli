package tags

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a tag by ID",
		RunE:  executeGet,
	}

	cmd.Flags().String("id", "", "Tag ID.")

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Tag:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		slog.Error("Failed to resolve tag ID", slog.Any("error", err))
		return err
	}

	tag, err := client.Tags.Get(cmd.Context(), id)
	if err != nil {
		slog.Error("Failed to get tag", slog.Any("error", err))
		return err
	}

	renderTagsTable([]linkwarden.TagInfo{*tag})

	return nil
}
