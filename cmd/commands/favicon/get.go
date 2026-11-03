package favicon

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Fetch a favicon for a URL",
		RunE:  executeGet,
	}

	cmd.Flags().String("url", "", "URL to fetch a favicon for.")
	cmd.Flags().String("output", "", "Path to save the favicon to. Without it, only content-type/size are printed.")

	if err := cmd.MarkFlagRequired("url"); err != nil {
		return nil
	}

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	targetURL, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	client := utils.NewLinkwardenClient()

	favicon, err := client.Favicon.Get(cmd.Context(), targetURL)
	if err != nil {
		slog.Error("Failed to fetch favicon", slog.Any("error", err))
		return err
	}

	if output == "" {
		fmt.Printf("Content-Type: %s, Size: %d bytes\n", favicon.ContentType, len(favicon.Data))
		return nil
	}

	if err := os.WriteFile(output, favicon.Data, 0o600); err != nil {
		slog.Error("Failed to write favicon file", slog.Any("error", err))
		return err
	}

	slog.Info("Favicon saved successfully", slog.String("path", output))

	return nil
}
