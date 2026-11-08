package tokens

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Revoke() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke a token",
		RunE:  executeRevoke,
	}

	utils.AddYesFlag(cmd)

	return cmd
}

func executeRevoke(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching tokens…")
	tokenList, err := client.Tokens.List(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get tokens: %w", err)
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

	ok, err := utils.Confirm(cmd, fmt.Sprintf("Revoke token %q? Anything using it will stop working.", result))
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was revoked")
		return nil
	}

	spinner = utils.StartSpinner("Revoking token…")
	err = client.Tokens.Revoke(cmd.Context(), tokenList[idx].ID)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}

	log.Info().Str("name", result).Msg("Token revoked successfully")

	return nil
}
