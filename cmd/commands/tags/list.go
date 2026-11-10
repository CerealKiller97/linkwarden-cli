package tags

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tags",
		Example: `  # List your tags
  linkwarden tags list

  # Tags that match a search
  linkwarden tags list --search go

  # Tag IDs only, one per line
  linkwarden tags list -q`,
		RunE: executeList,
	}

	cmd.Flags().String("search", "", "Search query string.")
	cmd.Flags().Int("cursor", 0, "Pagination cursor.")

	return cmd
}

func executeList(cmd *cobra.Command, _ []string) error {
	search, err := cmd.Flags().GetString("search")
	if err != nil {
		return err
	}
	cursor, err := cmd.Flags().GetInt("cursor")
	if err != nil {
		return err
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching tags…")
	tagList, err := client.Tags.List(cmd.Context(), linkwarden.ListTagsOptions{
		Search: search,
		Cursor: cursor,
	})
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("list tags: %w", err)
	}

	renderTagsTable(tagList, tagList)

	return nil
}
