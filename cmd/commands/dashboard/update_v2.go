package dashboard

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func UpdateV2() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-v2",
		Short: "Update dashboard layout (v2)",
		RunE:  executeUpdateV2,
	}

	cmd.Flags().String("file", "", "Path to a JSON file containing the layout, "+
		`e.g. [{"type":"STATS","enabled":true,"order":0}].`)

	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil
	}

	return cmd
}

func executeUpdateV2(cmd *cobra.Command, _ []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("Failed to read layout file", slog.Any("error", err))
		return err
	}

	var layout []linkwarden.DashboardLayoutItem
	if err := json.Unmarshal(fileBytes, &layout); err != nil {
		slog.Error("Failed to parse layout file", slog.Any("error", err))
		return err
	}

	client := utils.NewLinkwardenClient()

	data, err := client.Dashboard.UpdateV2Layout(cmd.Context(), layout)
	if err != nil {
		slog.Error("Failed to update dashboard layout", slog.Any("error", err))
		return err
	}

	renderDashboardV2Stats(data)

	return nil
}
