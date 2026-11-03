package search

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search links",
		RunE:  execute,
	}

	cmd.Flags().String("query", "", "Search query string.")
	cmd.Flags().Int("sort", 0, "Sort order.")
	cmd.Flags().Int("cursor", 0, "Pagination cursor.")

	return cmd
}

func execute(cmd *cobra.Command, _ []string) error {
	query, err := cmd.Flags().GetString("query")
	if err != nil {
		return err
	}
	cursor, err := cmd.Flags().GetInt("cursor")
	if err != nil {
		return err
	}

	var sort *int
	if cmd.Flags().Changed("sort") {
		s, err := cmd.Flags().GetInt("sort")
		if err != nil {
			return err
		}
		sort = &s
	}

	client := utils.NewLinkwardenClient()

	links, err := client.Search.Links(cmd.Context(), linkwarden.SearchOptions{
		SearchQueryString: query,
		Sort:              sort,
		Cursor:            cursor,
	})
	if err != nil {
		slog.Error("Failed to search links", slog.Any("error", err))
		return err
	}

	renderSearchResultsTable(links)

	return nil
}
