package tags

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Get(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a tag by ID",
		Example: `  # Pick the tag from a list
  linkwarden tags get

  # Get tag 14
  linkwarden tags get --id 14`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGet(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Tag ID.")

	return cmd
}

func executeGet(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		app, cmd,
		"id",
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

	renderTagsTable(app, []linkwarden.TagInfo{*tag}, tag)

	return nil
}
