package tags

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage tags",
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Create())
	cmd.AddCommand(Update())
	cmd.AddCommand(Delete())
	cmd.AddCommand(BulkDelete())

	return cmd
}
