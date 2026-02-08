package tokens

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/CerealKiller97/linkwarden-cli/pkg/services/token"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Expiry int

const (
	SevenDays  Expiry = 0
	ThirtyDays Expiry = 1
	SixtyDays  Expiry = 2
	NinetyDays Expiry = 3
	NoExpiry   Expiry = 4
)

func Generate() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "generate",
		RunE: execute,
	}

	return cmd
}

func execute(cmd *cobra.Command, _ []string) error {
	prompt := promptui.Prompt{
		Label: "Token name",
		Validate: func(s string) error {
			if len(s) == 0 {
				return errors.New("Token name is required")
			}

			return nil
		},
		Templates: nil,
		IsConfirm: false,
		IsVimMode: false,
	}

	promptResult, err := prompt.Run()
	if err != nil {
		log.Err(err).Msg("tokens name prompt failed")
		return err
	}

	expiration := promptui.Select{
		Label: "Token Expiration",
		Items: []string{"7 days", "30 days", "60 days", "90 days", "no expiration"},
	}

	_, expirationResult, err := expiration.Run()
	if err != nil {
		log.Err(err).Msg("tokens expiration prompt failed")
		return err
	}

	tokenService := token.New(
		utils.GetInstanceURL(),
		utils.GetLinkwardenToken(),
		zerolog.New(os.Stdout),
	)

	token, err := tokenService.Create(cmd.Context(), promptResult, int(stringToInt(expirationResult)))
	if err != nil {
		log.Err(err).Msg("tokens create prompt failed")

		return err
	}

	fmt.Println(os.Getenv("SHELL"))

	log.Info().Str("token", token).Msg("Token created")

	// Prompt do you want app to set ENV Variable on your system?
	// Only in Linux / macOS
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		prompt := promptui.Select{
			Label: "Do you want CLI to set token?",
			Items: []string{"Yes", "No"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			log.Err(err).Msg("tokens prompt failed")
			return err
		}

		if result == "Yes" {
			log.Info().Str("token", token).Msg("Token set")
			if err := utils.AppendTokenToShellFile(token); err != nil {
				log.Err(err).Msg("Failed to append token to shell file")
				return err
			}

			if err := utils.Reload(); err != nil {
				log.Err(err).Msg("Failed to reload shell file")
				return err
			}
		}
	}

	return nil
}

func stringToInt(tokenExpiration string) Expiry {
	expiration := NoExpiry
	switch tokenExpiration {
	case "no expiration":
		expiration = NoExpiry
	case "30 days":
		expiration = ThirtyDays
	case "60 days":
		expiration = SixtyDays
	case "90 days":
		expiration = NinetyDays
	case "7 days":
		expiration = SevenDays
	default:
		log.Info().Msgf("Invalid token expiration: %s", tokenExpiration)
	}

	return expiration
}
