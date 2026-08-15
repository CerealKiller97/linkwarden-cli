package rss

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rss",
		Short: "Manage RSS subscriptions",
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Create())
	cmd.AddCommand(Delete())

	return cmd
}
