package tokens

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func List(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available tokens",
		Example: `  # List your API tokens
  linkwarden tokens list

  # Same, as JSON
  linkwarden tokens list --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeList(app, cmd, args)
		},
	}
}

func executeList(app *app.App, cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching tokens…")
	tokenList, err := client.Tokens.List(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get tokens: %w", err)
	}

	t := utils.NewTable(app, "Tokens")
	t.Header("ID", "Name", "Session", "Expires", "Created")

	for _, token := range tokenList {
		t.Row(token.ID, token.Name, token.IsSession, token.Expires, token.CreatedAt)
	}

	t.JSON(tokenList)
	t.Render()

	return nil
}
