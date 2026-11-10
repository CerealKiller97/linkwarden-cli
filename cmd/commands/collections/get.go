package collections

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a collection by ID",
		RunE:  executeGet,
	}

	cmd.Flags().Int("id", 0, "Collection ID.")

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd,
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

	renderCollectionsTable([]linkwarden.Collection{*collection}, collection)

	return nil
}
