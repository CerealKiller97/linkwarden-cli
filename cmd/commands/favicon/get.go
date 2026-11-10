package favicon

import (
	"fmt"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Get(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Fetch a favicon for a URL",
		Example: `  # Print the favicon's type and size
  linkwarden favicon get --url https://go.dev

  # Save it to a file
  linkwarden favicon get --url https://go.dev --output go.ico`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeGet(app, cmd, args)
		},
	}

	cmd.Flags().String("url", "", "URL to fetch a favicon for.")
	cmd.Flags().String("output", "", "Path to save the favicon to. Without it, only content-type/size are printed.")

	if err := cmd.MarkFlagRequired("url"); err != nil {
		return nil
	}

	return cmd
}

func executeGet(app *app.App, cmd *cobra.Command, _ []string) error {
	targetURL, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching favicon…")
	favicon, err := client.Favicon.Get(cmd.Context(), targetURL)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("fetch favicon: %w", err)
	}

	if output == "" {
		return utils.PrintFileInfo(app, favicon.ContentType, len(favicon.Data))
	}

	if err := os.WriteFile(output, favicon.Data, 0o600); err != nil {
		return fmt.Errorf("write favicon file: %w", err)
	}

	log.Info().Str("path", output).Msg("Favicon saved successfully")

	return nil
}
