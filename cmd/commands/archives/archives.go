package archives

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "archives",
		Short: "Manage link archive files",
	}

	cmd.AddCommand(Get())
	cmd.AddCommand(Upload())
	cmd.AddCommand(UploadStandalone())

	return cmd
}
