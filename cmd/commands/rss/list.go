package rss

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List RSS subscriptions",
		RunE:  executeList,
	}

	return cmd
}

func executeList(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	subs, err := client.RSS.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("list RSS subscriptions: %w", err)
	}

	renderRSSTable(subs)

	return nil
}
