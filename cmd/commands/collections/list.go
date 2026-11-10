package collections

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List collections",
		Example: `  # List your collections
  linkwarden collections list

  # ID and name of each, tab-separated
  linkwarden collections list --json | jq -r '.[] | "\(.id)\t\(.name)"'`,
		RunE: executeList,
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

	renderCollectionsTable(collections, collections)

	// The list already holds the Unorganized collection, so refresh this
	// instance's cached ID for "collections create" at no extra cost. A cache
	// that can't be written isn't worth failing the list over.
	if id := unorganizedID(collections); id != 0 {
		if err := persistUnorganizedCollectionID(utils.GetInstanceURL(), id); err != nil {
			log.Warn().Err(err).Msg("Failed to cache the Unorganized collection ID")
		}
	}

	return nil
}
