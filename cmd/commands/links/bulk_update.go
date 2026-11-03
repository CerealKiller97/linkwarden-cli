package links

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func BulkUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-update",
		Short: "Apply a collection/tag change to multiple links",
		RunE:  executeBulkUpdate,
	}

	cmd.Flags().String("ids", "", "Comma-separated list of link IDs to update.")
	cmd.Flags().Int("collection-id", 0, "Move links to this collection ID.")
	cmd.Flags().StringSlice("tag", nil, "Tag name to apply (repeatable).")
	cmd.Flags().Bool("remove-previous-tags", false, "Remove tags not in --tag before applying.")

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeBulkUpdate(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		slog.Error("Failed to parse --ids", slog.Any("error", err))
		return err
	}

	collectionID, err := cmd.Flags().GetInt("collection-id")
	if err != nil {
		return err
	}
	tagNames, err := cmd.Flags().GetStringSlice("tag")
	if err != nil {
		return err
	}
	removePreviousTags, err := cmd.Flags().GetBool("remove-previous-tags")
	if err != nil {
		return err
	}

	links := make([]linkwarden.Link, len(ids))
	for i, id := range ids {
		links[i] = linkwarden.Link{ID: id}
	}

	tags := make([]linkwarden.Tag, len(tagNames))
	for i, name := range tagNames {
		tags[i] = linkwarden.Tag{Name: name}
	}

	client := utils.NewLinkwardenClient()

	err = client.Links.BulkUpdate(cmd.Context(), linkwarden.BulkUpdateLinksRequest{
		Links:              links,
		RemovePreviousTags: removePreviousTags,
		NewData: linkwarden.BulkUpdateLinksNewData{
			CollectionID: collectionID,
			Tags:         tags,
		},
	})
	if err != nil {
		slog.Error("Failed to bulk update links", slog.Any("error", err))
		return err
	}

	slog.Info("Links updated successfully", slog.Any("ids", ids))

	return nil
}
