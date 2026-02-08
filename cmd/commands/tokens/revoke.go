package tokens

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/services/token"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
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
	tokenService := token.New(
		utils.GetInstanceURL(),
		utils.GetLinkwardenToken(),
		zerolog.New(os.Stdout),
	)

	tokens, err := tokenService.Get(cmd.Context())
	if err != nil {
		log.Err(err).Msg("Error getting tokens")
		return err
	}

	tokenNames := make([]string, len(tokens))
	tokenMap := make(map[string]int, len(tokens))
	for i, token := range tokens {
		tokenNames[i] = token.Name
		tokenMap[token.Name] = token.ID
	}

	prompt := promptui.Select{
		Label: "Select Token",
		Items: tokenNames,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	if err := tokenService.Revoke(cmd.Context(), tokenMap[result]); err != nil {
		log.Err(err).Msg("Error revoking token")
		return err
	}

	log.Info().Str("name", result).Msg("Token revoked successfully")

	return nil
}
