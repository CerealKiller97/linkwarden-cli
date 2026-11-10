package search

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/choice"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

// sortOrders are Linkwarden's link sort orders, in the order of the API's codes.
var sortOrders = []choice.Option[int]{
	{Name: "newest", Value: 0, Description: "Date added, newest first"},
	{Name: "oldest", Value: 1, Description: "Date added, oldest first"},
	{Name: "name-asc", Value: 2, Description: "Name, A to Z"},
	{Name: "name-desc", Value: 3, Description: "Name, Z to A"},
	{Name: "description-asc", Value: 4, Description: "Description, A to Z"},
	{Name: "description-desc", Value: 5, Description: "Description, Z to A"},
}

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search links",
		Example: `  # Search your links
  linkwarden search --query "kubernetes operators"

  # Oldest matches first
  linkwarden search --query golang --sort oldest

  # IDs only, one per line
  linkwarden search --query golang -q`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(app, cmd, args)
		},
	}

	cmd.Flags().String("query", "", "Search query string.")
	choice.Flag(cmd, "sort", "", "Sort order (the server's default when omitted).", sortOrders)
	cmd.Flags().Int("cursor", 0, "Pagination cursor.")

	return cmd
}

func execute(app *app.App, cmd *cobra.Command, _ []string) error {
	query, err := cmd.Flags().GetString("query")
	if err != nil {
		return err
	}
	cursor, err := cmd.Flags().GetInt("cursor")
	if err != nil {
		return err
	}

	var sort *int
	s, set, err := choice.Get[int](cmd, "sort")
	if err != nil {
		return err
	}
	if set {
		sort = &s
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Searching links…")
	links, err := client.Search.Links(cmd.Context(), linkwarden.SearchOptions{
		SearchQueryString: query,
		Sort:              sort,
		Cursor:            cursor,
	})
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("search links: %w", err)
	}

	renderSearchResultsTable(app, links)

	return nil
}
