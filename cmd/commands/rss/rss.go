package rss

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/spf13/cobra"
)

func Command(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rss",
		Short: "Manage RSS subscriptions",
	}

	cmd.AddCommand(List(app))
	cmd.AddCommand(Create(app))
	cmd.AddCommand(Delete(app))

	return cmd
}
