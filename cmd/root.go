package cmd

import (
	"errors"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/archives"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/collections"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/config"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/dashboard"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/favicon"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/highlights"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/migrations"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/rss"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/search"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/links"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/tags"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/tokens"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func registerCommands(root *cobra.Command) {
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
		utils.ConfigureDefaultLogger("info")

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")
		if err != nil {
			return err
		}
		utils.SetNonInteractive(nonInteractive)

		linkwardenToken, exists := os.LookupEnv("LINKWARDEN_TOKEN")

		if !exists {
			return errors.New("environment variable LINKWARDEN_TOKEN is not set")
		}

		// TODO: Trim / if present from URL
		instanceURL, exists := os.LookupEnv("LINKWARDEN_URL")
		if !exists {
			instanceURL = "https://cloud.linkwarden.app"
		}

		utils.SetVars(linkwardenToken, instanceURL)

		// INSTANCE_URL env var if not set use https://cloud.linkwarden.app/

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

// version is the released version, set at build time by GoReleaser (.goreleaser.yaml) via
// -ldflags "-X github.com/CerealKiller97/linkwarden-cli/cmd.version=1.2.3".
var version = "dev"

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "linkwarden",
		Short:   "Linkwarden is a CLI and SDK for the Linkwarden platform",
		Long:    `Linkwarden is a CLI and SDK for the Linkwarden platform`,
		Version: version,
	}

	registerCommands(rootCmd)

	return rootCmd
}

func Execute() {
	if err := Root().Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute CLI command")
	}
}
