package rss

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Msg("Failed to list RSS subscriptions")
		return err
	}

	renderRSSTable(subs)

	return nil
}
