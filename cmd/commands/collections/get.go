package collections

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
		Short: "Get a collection by ID",
		Example: `  # Pick the collection from a list
  linkwarden collections get

  # Get collection 7
  linkwarden collections get --id 7

  # Print only its name
  linkwarden collections get --id 7 --json | jq -r .name`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGet(app, cmd, args)
		},
	}

	cmd.Flags().Int("id", 0, "Collection ID.")

	return cmd
}

func executeGet(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		app, cmd,
		"id",
		"Select Collection:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve collection ID: %w", err)
	}

	spinner := utils.StartSpinner("Fetching collection…")
	collection, err := client.Collections.Get(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get collection: %w", err)
	}

	renderCollectionsTable(app, []linkwarden.Collection{*collection}, collection)

	return nil
}
