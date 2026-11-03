package rss

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an RSS subscription",
		RunE:  executeCreate,
	}

	cmd.Flags().String("name", "", "Subscription name.")
	cmd.Flags().String("url", "", "Feed URL.")
	cmd.Flags().Int("collection-id", 0, "Destination collection ID.")
	cmd.Flags().String("collection-name", "", "Destination collection name (used if collection-id is not set).")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}
	if err := cmd.MarkFlagRequired("url"); err != nil {
		return nil
	}

	return cmd
}

func executeCreate(cmd *cobra.Command, _ []string) error {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}
	collectionID, err := cmd.Flags().GetInt("collection-id")
	if err != nil {
		return err
	}
	collectionName, err := cmd.Flags().GetString("collection-name")
	if err != nil {
		return err
	}

	client := utils.NewLinkwardenClient()

	sub, err := client.RSS.Create(cmd.Context(), linkwarden.CreateRSSRequest{
		Name:           name,
		URL:            url,
		CollectionID:   collectionID,
		CollectionName: collectionName,
	})
	if err != nil {
		slog.Error("Failed to create RSS subscription", slog.Any("error", err))
		return err
	}

	renderRSSTable([]linkwarden.RSSSubscription{*sub})

	return nil
}
