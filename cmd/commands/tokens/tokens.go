package tokens

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Manage tokens",
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Generate())
	cmd.AddCommand(Revoke())

	return cmd
}
