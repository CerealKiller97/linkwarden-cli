package collections

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
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
		"--id",
		"Select Collection:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		slog.Error("Failed to resolve collection ID", slog.Any("error", err))
		return err
	}

	collection, err := client.Collections.Get(cmd.Context(), id)
	if err != nil {
		slog.Error("Failed to get collection", slog.Any("error", err))
		return err
	}

	renderCollectionsTable([]linkwarden.Collection{*collection})

	return nil
}
