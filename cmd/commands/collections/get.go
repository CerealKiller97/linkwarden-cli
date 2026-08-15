package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a collection by ID",
		RunE:  executeGet,
	}

	cmd.Flags().String("id", "", "Collection ID.")

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"Select Collection:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve collection ID")
		return err
	}

	collection, err := client.Collections.Get(cmd.Context(), id)
	if err != nil {
		log.Err(err).Msg("Failed to get collection")
		return err
	}

	renderCollectionsTable([]linkwarden.Collection{*collection})

	return nil
}
