package links

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "links",
		Short: "Manage links",
	}

	cmd.AddCommand(Create())
	cmd.AddCommand(Import())
	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Update())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Archive())
	cmd.AddCommand(BulkUpdate())
	cmd.AddCommand(BulkDelete())
	cmd.AddCommand(DeleteArchive())
	cmd.AddCommand(GetHighlights())

	return cmd
}
