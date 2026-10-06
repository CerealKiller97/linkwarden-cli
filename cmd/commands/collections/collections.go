package collections

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collections",
		Short: "Manage collections",
	}

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Update())
	cmd.AddCommand(Delete())

	return cmd
}
