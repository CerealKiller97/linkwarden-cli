package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/internal/pickers"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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

	id, err := utils.ResolveIDPaged(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Link:",
		pickers.Links(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve link ID: %w", err)
	}

	spinner := utils.StartSpinner("Fetching link…")
	link, err := client.Links.Get(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get link: %w", err)
	}

	renderLinkTable(link)

	return nil
}
