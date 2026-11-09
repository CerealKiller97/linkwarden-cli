package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
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

	cobra.CheckErr(cmd.MarkFlagRequired("ids"))

	return cmd
}

func executeBulkUpdate(cmd *cobra.Command, _ []string) error {
	ids, err := parseIDs(cmd.Flag("ids").Value.String())
	if err != nil {
		return fmt.Errorf("parse --ids: %w", err)
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

	spinner := utils.StartSpinner("Updating links…")
	err = client.Links.BulkUpdate(cmd.Context(), linkwarden.BulkUpdateLinksRequest{
		Links:              links,
		RemovePreviousTags: removePreviousTags,
		NewData: linkwarden.BulkUpdateLinksNewData{
			CollectionID: collectionID,
			Tags:         tags,
		},
	})
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("bulk update links: %w", err)
	}

	log.Info().Ints("ids", ids).Msg("Links updated successfully")

	return nil
}
