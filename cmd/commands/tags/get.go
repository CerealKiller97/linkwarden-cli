package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Get() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a tag by ID",
		RunE:  executeGet,
	}

	cmd.Flags().String("id", "", "Tag ID.")

	return cmd
}

func executeGet(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Tag:",
		tagsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve tag ID")
		return err
	}

	tag, err := client.Tags.Get(cmd.Context(), id)
	if err != nil {
		log.Err(err).Msg("Failed to get tag")
		return err
	}

	renderTagsTable([]linkwarden.TagInfo{*tag})

	return nil
}
