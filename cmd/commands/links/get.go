package links

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a link by ID",
		RunE:  executeGet,
	}

	cmd.Flags().String("id", "", "Link ID.")

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Link:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		slog.Error("Failed to resolve link ID", slog.Any("error", err))
		return err
	}

	link, err := client.Links.Get(cmd.Context(), id)
	if err != nil {
		slog.Error("Failed to get link", slog.Any("error", err))
		return err
	}

	renderLinksTable([]linkwarden.Link{*link})

	return nil
}
