package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/archives"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/collections"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/config"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/dashboard"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/favicon"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/highlights"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/initcmd"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/migrations"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/rss"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/search"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/links"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/tags"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/tokens"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
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

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Flags and args have parsed by now, so any error from here on is a runtime
		// failure, not a usage mistake: don't bury it under the help text.
		cmd.SilenceUsage = true

		utils.ConfigureDefaultLogger("info")

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")
		if err != nil {
			return err
		}
		utils.SetNonInteractive(nonInteractive)

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

		// TODO: App set
		// configure logger
		// initialize DI container

		return nil
	}

	//	appLogger.ConfigureDefaultLogger(cfg.Logging.Level, cfg.Logging.PrettyPrint)
	//	//nolint:all
	//	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	//	command.SetContext(context.WithValue(ctx, constants.ContextKeyContainer, container.New(ctx, cfg)))
	//	return nil
	//})

	//root.PersistentPostRunE = cmd.PersistentPostRunE(func(cmd *cobra.Command, _ []string) error {
	//	return cmd.Context().Value(constants.ContextKeyContainer).(*container.Container).Close()
	//})
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
	if err := Root().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
