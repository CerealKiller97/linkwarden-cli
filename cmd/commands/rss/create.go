package rss

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an RSS subscription",
		Example: `  # Follow a feed into collection 3
  linkwarden rss create --name "The Go Blog" --url https://go.dev/blog/feed.atom --collection-id 3

  # Name the destination collection instead
  linkwarden rss create --name "The Go Blog" --url https://go.dev/blog/feed.atom --collection-name Go`,
		RunE: executeCreate,
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

	spinner := utils.StartSpinner("Adding RSS subscription…")
	sub, err := client.RSS.Create(cmd.Context(), linkwarden.CreateRSSRequest{
		Name:           name,
		URL:            url,
		CollectionID:   collectionID,
		CollectionName: collectionName,
	})
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("create RSS subscription: %w", err)
	}

	renderRSSTable(cmd.Context(), []linkwarden.RSSSubscription{*sub}, sub)

	return nil
}
