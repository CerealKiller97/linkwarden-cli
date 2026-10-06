package favicon

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "favicon",
		Short: "Fetch favicons",
	}

	cmd.AddCommand(Get())

	return cmd
}
