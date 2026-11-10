package tokens

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available tokens",
		Example: `  # List your API tokens
  linkwarden tokens list

  # Same, as JSON
  linkwarden tokens list --json`,
		RunE: executeList,
	}
}

func executeList(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching tokens…")
	tokenList, err := client.Tokens.List(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get tokens: %w", err)
	}

	t := utils.NewTable(cmd.Context(), "Tokens")
	t.Header("ID", "Name", "Session", "Expires", "Created")

	for _, token := range tokenList {
		t.Row(token.ID, token.Name, token.IsSession, token.Expires, token.CreatedAt)
	}

	t.JSON(tokenList)
	t.Render()

	return nil
}
