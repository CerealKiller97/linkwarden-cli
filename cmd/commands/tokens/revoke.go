package tokens

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Revoke() *cobra.Command {
	return &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a token",
		RunE:  executeRevoke,
	}
}

func executeRevoke(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	tokenList, err := client.Tokens.List(cmd.Context())
	if err != nil {
		log.Err(err).Msg("Error getting tokens")
		return err
	}

	tokenNames := make([]string, len(tokenList))
	for i, token := range tokenList {
		tokenNames[i] = token.Name
	}

	// Revoke by index rather than by name: two tokens can share a name.
	idx, result, err := utils.Select("Select token to revoke:", tokenNames)
	if err != nil {
		return err
	}

	if err := client.Tokens.Revoke(cmd.Context(), tokenList[idx].ID); err != nil {
		log.Err(err).Msg("Error revoking token")
		return err
	}

	log.Info().Str("name", result).Msg("Token revoked successfully")

	return nil
}
