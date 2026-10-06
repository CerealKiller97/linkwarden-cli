package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Msg("Failed to list collections")
		return err
	}

	renderCollectionsTable(collections)

	return nil
}
