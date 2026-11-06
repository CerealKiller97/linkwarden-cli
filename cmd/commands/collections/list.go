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

	collections, err := client.Collections.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("list collections: %w", err)
	}

	renderCollectionsTable(collections)

	return nil
}
