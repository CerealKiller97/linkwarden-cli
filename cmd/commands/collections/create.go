package collections

import "github.com/spf13/cobra"

func Create() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a collection",
		Long:  `Create a collection`,
		RunE:  executeCreate,
	}
}

func executeCreate(cmd *cobra.Command, args []string) error {
	return nil
}
