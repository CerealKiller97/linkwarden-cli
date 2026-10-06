package tokens

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available tokens",
		RunE:  executeList,
	}
}

func executeList(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	tokenList, err := client.Tokens.List(cmd.Context())
	if err != nil {
		log.Err(err).Msg("Error getting tokens")
		return err
	}

	t := utils.NewTable("Tokens")
	t.AppendHeader(table.Row{"ID", "Name", "Is Session", "Expires", "Created At"})

	for _, token := range tokenList {
		t.AppendRow(table.Row{token.ID, token.Name, token.IsSession, token.Expires, token.CreatedAt})
	}

	t.AppendSeparator()
	t.Render()

	return nil
}
