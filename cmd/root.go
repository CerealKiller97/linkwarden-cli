package cmd

import (
	"errors"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/collections"
	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/migrations"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/cmd/commands/links"
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

	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		utils.ConfigureDefaultLogger("info")

		linkwardenToken, exists := os.LookupEnv("LINKWARDEN_TOKEN")

		if !exists {
			return errors.New("Environment variable LINKWARDEN_TOKEN is not set")
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

func Root() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "linkwarden-cli",
		Short: "LinkwardenCLI is CLI SDK for the Linkwarden platform",
		Long:  `LinkwardenCLI is CLI SDK for the Linkwarden platform`,
	}

	registerCommands(rootCmd)

	return rootCmd
}

func Execute() {
	if err := Root().Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute CLI command")
	}
}
