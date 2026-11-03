package collections

import (
	"log/slog"

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
		slog.Error("Failed to list collections", slog.Any("error", err))
		return err
	}

	renderCollectionsTable(collections)

	return nil
}
