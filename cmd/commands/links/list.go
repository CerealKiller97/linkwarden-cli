package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List links",
		Example: `  # Your newest links
  linkwarden links list

  # Links in collection 3
  linkwarden links list --collection-id 3

  # Pinned links with tag 14
  linkwarden links list --tag-id 14 --pinned-only

  # Links that match a search
  linkwarden links list --search kubernetes

  # The URL of each link, one per line
  linkwarden links list --json | jq -r '.[].url'`,
		RunE: executeList,
	}

	cmd.Flags().Int("collection-id", 0, "Filter by collection ID.")
	cmd.Flags().Int("tag-id", 0, "Filter by tag ID.")
	cmd.Flags().Bool("pinned-only", false, "Only show pinned links.")
	cmd.Flags().String("search", "", "Search query string.")
	cmd.Flags().Int("cursor", 0, "Pagination cursor.")

	return cmd
}

func executeList(cmd *cobra.Command, _ []string) error {
	collectionID, err := cmd.Flags().GetInt("collection-id")
	if err != nil {
		return err
	}
	tagID, err := cmd.Flags().GetInt("tag-id")
	if err != nil {
		return err
	}
	pinnedOnly, err := cmd.Flags().GetBool("pinned-only")
	if err != nil {
		return err
	}
	search, err := cmd.Flags().GetString("search")
	if err != nil {
		return err
	}
	cursor, err := cmd.Flags().GetInt("cursor")
	if err != nil {
		return err
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching links…")
	links, err := client.Links.List(cmd.Context(), linkwarden.ListLinksOptions{
		Cursor:            cursor,
		CollectionID:      collectionID,
		TagID:             tagID,
		PinnedOnly:        pinnedOnly,
		SearchQueryString: search,
	})
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("list links: %w", err)
	}

	renderLinksTable(cmd.Context(), links)

	return nil
}
