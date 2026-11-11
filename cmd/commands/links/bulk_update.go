package links

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func BulkUpdate(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulk-update",
		Short: "Apply a collection/tag change to multiple links",
		Example: `  # Move three links into collection 7
  linkwarden links bulk-update --ids 412,418,421 --collection-id 7

  # Replace their tags
  linkwarden links bulk-update --ids 412,418,421 --tag reading --tag weekend --remove-previous-tags`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeBulkUpdate(app, cmd, args)
		},
	}

	utils.AddIDsFlag(cmd, "Link IDs to update.")
	cmd.Flags().Int("collection-id", 0, "Move links to this collection ID.")
	cmd.Flags().StringSlice("tag", nil, "Tag name to apply (repeatable).")
	cmd.Flags().Bool("remove-previous-tags", false, "Remove tags not in --tag before applying.")

	if err := cmd.MarkFlagRequired("ids"); err != nil {
		return nil
	}

	return cmd
}

func executeBulkUpdate(app *app.App, cmd *cobra.Command, _ []string) error {
	ids, err := utils.GetIDs(cmd)
	if err != nil {
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
