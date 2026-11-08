package collections

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List collections",
		RunE:  executeList,
	}
}

func executeList(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching collections…")
	collections, err := client.Collections.List(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("list collections: %w", err)
	}

	renderCollectionsTable(collections)

	return nil
}
