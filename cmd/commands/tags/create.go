package tags

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create or update tags in bulk",
		RunE:  executeCreate,
	}

	cmd.Flags().StringSlice("label", nil, "Tag label (repeatable).")

	if err := cmd.MarkFlagRequired("label"); err != nil {
		return nil
	}

	return cmd
}

func executeCreate(cmd *cobra.Command, _ []string) error {
	labels, err := cmd.Flags().GetStringSlice("label")
	if err != nil {
		return err
	}

	inputs := make([]linkwarden.CreateOrUpdateTagInput, len(labels))
	for i, label := range labels {
		inputs[i] = linkwarden.CreateOrUpdateTagInput{Label: label}
	}

	client := utils.NewLinkwardenClient()

	createdTags, err := client.Tags.CreateOrUpdateBulk(cmd.Context(), inputs)
	if err != nil {
		log.Err(err).Msg("Failed to create tags")
		return err
	}

	renderTagsTable(createdTags)

	return nil
}
