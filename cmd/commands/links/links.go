package links

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "Manage links",
	}

	cmd.AddCommand(Create())
	return cmd
}
