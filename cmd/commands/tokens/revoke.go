package tokens

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/manifoldco/promptui"
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
		slog.Error("Error getting tokens", slog.Any("error", err))
		return err
	}

	tokenNames := make([]string, len(tokenList))
	tokenMap := make(map[string]int, len(tokenList))
	for i, token := range tokenList {
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

	if err := client.Tokens.Revoke(cmd.Context(), tokenMap[result]); err != nil {
		slog.Error("Error revoking token", slog.Any("error", err))
		return err
	}

	slog.Info("Token revoked successfully", slog.String("name", result))

	return nil
}
