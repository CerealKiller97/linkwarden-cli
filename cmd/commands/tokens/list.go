package tokens

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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
	t.Columns(
		utils.Column{Header: "ID"},
		utils.Column{Header: "Name", Flex: true},
		utils.Column{Header: "Is Session"},
		utils.Column{Header: "Expires"},
		utils.Column{Header: "Created At"},
	)

	for _, token := range tokenList {
		t.Row(token.ID, token.Name, token.IsSession, token.Expires, token.CreatedAt)
	}

	t.Render()

	return nil
}
