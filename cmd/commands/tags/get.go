package tags

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
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
		return fmt.Errorf("resolve tag ID: %w", err)
	}

	spinner := utils.StartSpinner("Fetching tag…")
	tag, err := client.Tags.Get(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get tag: %w", err)
	}

	renderTagsTable([]linkwarden.TagInfo{*tag})

	return nil
}
