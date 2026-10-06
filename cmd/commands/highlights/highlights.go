package highlights

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "highlights",
		Short: "Manage link highlights",
	}

	cmd.AddCommand(Create())
	cmd.AddCommand(Delete())

	return cmd
}
