package config

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View instance configuration",
	}

	cmd.AddCommand(Get())

	return cmd
}
