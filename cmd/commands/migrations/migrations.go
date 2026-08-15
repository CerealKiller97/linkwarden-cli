package migrations

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrations",
		Short: "Manage migration (Import / Export)",
	}

	cmd.AddCommand(Import())
	cmd.AddCommand(Export())

	return cmd
}
