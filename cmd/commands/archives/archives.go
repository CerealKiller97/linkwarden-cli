package archives

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archives",
		Short: "Manage link archive files",
	}

	cmd.AddCommand(Get(app))
	cmd.AddCommand(Upload(app))
	cmd.AddCommand(UploadStandalone(app))

	return cmd
}
