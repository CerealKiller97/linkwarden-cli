package dashboard

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func UpdateV2(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-v2",
		Short: "Update dashboard layout (v2)",
		Example: `  # layout.json: [{"type":"STATS","enabled":true,"order":0},{"type":"RECENT_LINKS","enabled":true,"order":1}]
  linkwarden dashboard update-v2 --file layout.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUpdateV2(app, cmd, args)
		},
	}

	cmd.Flags().String("file", "", "Path to a JSON file containing the layout, "+
		`e.g. [{"type":"STATS","enabled":true,"order":0}].`)

	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil
	}

	return cmd
}

func executeUpdateV2(app *app.App, cmd *cobra.Command, _ []string) error {
	filePath, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read layout file: %w", err)
	}

	var layout []linkwarden.DashboardLayoutItem
	if err := json.Unmarshal(fileBytes, &layout); err != nil {
		return fmt.Errorf("parse layout file: %w", err)
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Updating dashboard layout…")
	data, err := client.Dashboard.UpdateV2Layout(cmd.Context(), layout)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("update dashboard layout: %w", err)
	}

	renderDashboardV2Stats(app, data)

	return nil
}
