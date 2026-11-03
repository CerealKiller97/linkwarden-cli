package links

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List links",
		RunE:  executeList,
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

	links, err := client.Links.List(cmd.Context(), linkwarden.ListLinksOptions{
		Cursor:            cursor,
		CollectionID:      collectionID,
		TagID:             tagID,
		PinnedOnly:        pinnedOnly,
		SearchQueryString: search,
	})
	if err != nil {
		slog.Error("Failed to list links", slog.Any("error", err))
		return err
	}

	renderLinksTable(links)

	return nil
}
