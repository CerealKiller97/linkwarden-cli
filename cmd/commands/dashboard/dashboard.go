package dashboard

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "View and configure the dashboard",
	}

	cmd.AddCommand(Get())
	cmd.AddCommand(GetV2())
	cmd.AddCommand(UpdateV2())

	return cmd
}
