package tokens

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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

func Generate(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Create an API token",
		Long: `Create an API token for another machine or script.

Prompts for a token name and an expiry (7, 30, 60 or 90 days, or never), then
prints the new token. On Linux and macOS it offers to add the token to your
shell config. Needs an interactive terminal.`,
		Example: `  # Asks for a name and an expiry, then prints the new token
  linkwarden tokens generate`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(app, cmd, args)
		},
	}

	return cmd
}

func execute(app *app.App, cmd *cobra.Command, _ []string) error {
	promptResult, err := utils.Input(app, "Token name:", func(s string) error {
		if strings.TrimSpace(s) == "" {
			return errors.New("token name is required")
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("prompt for token name: %w", err)
	}

	_, expirationResult, err := utils.Select(app, "Token expiration:",
		[]string{"7 days", "30 days", "60 days", "90 days", "no expiration"})
	if err != nil {
		return fmt.Errorf("prompt for token expiration: %w", err)
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Creating token…")
	created, err := client.Tokens.Create(cmd.Context(), promptResult, int(stringToInt(expirationResult)))
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("prompt for token confirmation: %w", err)
	}

	token := created.SecretKey

	log.Info().Str("token", token).Msg("Token created")

	// Prompt do you want app to set ENV Variable on your system?
	// Only in Linux / macOS
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		_, result, err := utils.Select(app, "Do you want the CLI to set the token in your shell config?",
			[]string{"Yes", "No"})
		if err != nil {
			return fmt.Errorf("prompt: %w", err)
		}

		if result == "Yes" {
			if err := utils.AppendTokenToShellFile(token); err != nil {
				return fmt.Errorf("append token to shell file: %w", err)
			}

			log.Info().
				Str("path", utils.ShellConfigPath()).
				Msg("Token added to shell config. Run `source` on it or restart your terminal to use it.")
		} else {
			printManualSetupInstructions(token)
		}
	}

	return nil
}

// printManualSetupInstructions prints copy-pasteable, colorized steps for
// exporting LINKWARDEN_TOKEN by hand, for users who decline to have the CLI
// edit their shell config directly. Color is skipped automatically when
// stdout isn't a terminal or NO_COLOR is set (see utils.ColorEnabled).
func printManualSetupInstructions(token string) {
	shell := filepath.Base(os.Getenv("SHELL"))

	var configFile, exportLine string
	switch shell {
	case "fish":
		configFile = utils.ShellConfigPath()
		exportLine = fmt.Sprintf("set -x LINKWARDEN_TOKEN %q", token)
	case "bash":
		configFile = "~/.bashrc"
		exportLine = fmt.Sprintf("export LINKWARDEN_TOKEN=%q", token)
	case "zsh":
		configFile = "~/.zshenv"
		exportLine = fmt.Sprintf("export LINKWARDEN_TOKEN=%q", token)
	default:
		configFile = "your shell's config file"
		exportLine = fmt.Sprintf("export LINKWARDEN_TOKEN=%q", token)
	}

	heading := lipgloss.NewStyle().Bold(true).Foreground(tui.Sky)
	step := heading
	cmdColor := lipgloss.NewStyle().Foreground(tui.Ocean)
	warn := lipgloss.NewStyle().Bold(true).Foreground(tui.Warning)

	on := utils.ColorEnabled()
	c := func(style lipgloss.Style, s string) string {
		if !on {
			return s
		}
		return style.Render(s)
	}

	printStep := func(n int, label, command string) {
		fmt.Printf("%s %s\n", c(step, fmt.Sprintf("%d.", n)), label)
		fmt.Printf("   %s\n", c(cmdColor, command))
	}

	fmt.Println()
	fmt.Println(c(heading, "No problem — the CLI won't touch your shell config. Set it up manually:"))
	fmt.Println()

	printStep(1, "For the current terminal session only, run:", exportLine)
	fmt.Println()
	printStep(2, fmt.Sprintf("To persist it across sessions, add that line to %s:", configFile),
		fmt.Sprintf("echo '%s' >> %s", exportLine, configFile))
	fmt.Println()
	printStep(3, "Reload your shell config (or restart your terminal):", "source "+configFile)
	fmt.Println()
	printStep(4, "Verify it's set:", "echo $LINKWARDEN_TOKEN")
	fmt.Println()

	fmt.Println(c(warn, "⚠ Keep this token secret — anyone with it can act as you on your Linkwarden instance."))
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
