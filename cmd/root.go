package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/archives"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/collections"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/config"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/dashboard"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/favicon"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/highlights"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/initcmd"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/links"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/migrations"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/rss"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/search"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/tags"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/tokens"
	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func registerCommands(root *cobra.Command) {
	root.AddCommand(initcmd.Command())
	root.AddCommand(links.Command())
	root.AddCommand(tokens.Command())
	root.AddCommand(collections.Command())
	root.AddCommand(migrations.Command())
	root.AddCommand(tags.Command())
	root.AddCommand(rss.Command())
	root.AddCommand(search.Command())
	root.AddCommand(dashboard.Command())
	root.AddCommand(highlights.Command())
	root.AddCommand(archives.Command())
	root.AddCommand(config.Command())
	root.AddCommand(favicon.Command())

	root.PersistentFlags().Bool("non-interactive", false, "Never prompt (e.g. for a picker); fail fast instead. "+
		"On by default when stdin isn't a terminal.")
	root.PersistentFlags().Bool("json", false, "Print results as JSON instead of a table.")
	root.PersistentFlags().BoolP("quiet", "q", false, "Print only the ID of each result, one per line, and hide info messages.")

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Flags and args have parsed by now, so any error from here on is a runtime
		// failure, not a usage mistake: don't bury it under the help text.
		cmd.SilenceUsage = true

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")
		if err != nil {
			return err
		}
		utils.SetNonInteractive(nonInteractive)

		if err := applyOutputFlags(cmd); err != nil {
			return err
		}

		if cmd.Annotations[initcmd.SkipAuthAnnotation] != "" {
			return nil
		}

		cfg, err := utils.LoadConfig()
		if err != nil {
			return fmt.Errorf("reading config file: %w", err)
		}

		token, instanceURL, err := resolveCredentials(os.Getenv("LINKWARDEN_TOKEN"), os.Getenv("LINKWARDEN_URL"), cfg)
		if err != nil {
			return err
		}

		utils.SetVars(token, instanceURL)

		return nil
	}
}

var errJSONAndQuiet = errors.New("--json and --quiet can't be used together")

// applyOutputFlags sets the output mode from --json and --quiet. --quiet also
// hides info messages on stderr, leaving warnings and errors.
func applyOutputFlags(cmd *cobra.Command) error {
	asJSON, err := cmd.Flags().GetBool("json")
	if err != nil {
		return err
	}
	quiet, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		return err
	}

	switch {
	case asJSON && quiet:
		return errJSONAndQuiet
	case asJSON:
		utils.SetOutputMode(utils.OutputJSON)
	case quiet:
		utils.SetOutputMode(utils.OutputIDs)
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		utils.SetOutputMode(utils.OutputTable)
	}

	return nil
}

// errNotLoggedIn is returned when neither LINKWARDEN_TOKEN nor the config file
// provides a token.
var errNotLoggedIn = errors.New(`not logged in: run "linkwarden init" or set LINKWARDEN_TOKEN`)

// resolveCredentials picks the token and instance URL to use. LINKWARDEN_TOKEN
// wins over the config file; with it, the URL comes from LINKWARDEN_URL or the
// default, as before the config file existed. Without it, the config file's token
// is used, with its URL unless LINKWARDEN_URL overrides it.
func resolveCredentials(envToken, envURL string, cfg utils.Config) (string, string, error) {
	if envToken != "" {
		return envToken, utils.NormalizeInstanceURL(envURL), nil
	}

	if cfg.Token == "" {
		return "", "", errNotLoggedIn
	}

	instanceURL := cfg.URL
	if envURL != "" {
		instanceURL = envURL
	}

	return cfg.Token, utils.NormalizeInstanceURL(instanceURL), nil
}

// version is the released version, set at build time by GoReleaser (.goreleaser.yaml) via
// -ldflags "-X github.com/CerealKiller97/linkwarden-cli/cmd.version=1.2.3".
var version = "dev"

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "linkwarden",
		Short:   "Linkwarden is a CLI for the Linkwarden platform",
		Long:    `Linkwarden is a CLI for the Linkwarden platform`,
		Version: version,
		// Execute prints the error itself, once.
		SilenceErrors: true,
	}

	registerCommands(rootCmd)

	return rootCmd
}

func Execute() {
	// Configured here rather than in PersistentPreRunE so flag and argument
	// errors, which fail before any PreRun hook, are logged the same way.
	utils.ConfigureDefaultLogger("info")

	ctx, stop := interruptContext()
	err := Root().ExecuteContext(ctx)
	stop()

	switch {
	case err == nil, errors.Is(err, tui.ErrAborted):
	case errors.Is(err, context.Canceled):
		// The wrapped error ends in "context canceled", which reads like a bug.
		log.Warn().Msg("Interrupted")
	default:
		log.Error().Msg(err.Error())
	}

	if code := exitCode(err); code != 0 {
		os.Exit(code)
	}
}

// interruptContext returns a context that is cancelled on Ctrl+C (SIGINT) or
// SIGTERM, so in-flight API requests stop and the command can return cleanly.
// After the first signal the default handling comes back, so a second Ctrl+C
// kills the process even if something isn't watching the context.
func interruptContext() (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop()
	}()

	return ctx, stop
}

// exitAborted is the exit status when the user cancels a prompt with Esc or
// Ctrl+C, or interrupts a running command: 128+SIGINT, as shells report an
// interrupted command.
const exitAborted = 130

// isInterrupted reports whether err means the user stopped the command: a
// cancelled prompt, or Ctrl+C cancelling the command's context.
func isInterrupted(err error) bool {
	return errors.Is(err, tui.ErrAborted) || errors.Is(err, context.Canceled)
}

// exitCode maps the error Execute got back to the process exit status. Stopping
// the command is the user's choice, not a failure, so it exits quietly with
// exitAborted rather than logging an error.
func exitCode(err error) int {
	switch {
	case err == nil:
		return 0
	case isInterrupted(err):
		return exitAborted
	default:
		return 1
	}
}
