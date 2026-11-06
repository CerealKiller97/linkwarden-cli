package links

import (
	"fmt"
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
		return fmt.Errorf("resolve link ID: %w", err)
	}

	link, err := client.Links.Get(cmd.Context(), id)
	if err != nil {
		return fmt.Errorf("get link: %w", err)
	}

	renderLinksTable([]linkwarden.Link{*link})

	return nil
}
