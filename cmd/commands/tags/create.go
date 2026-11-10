package tags

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/spf13/cobra"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create or update tags in bulk",
		Example: `  # Create a tag
  linkwarden tags create --label golang

  # Create several at once
  linkwarden tags create --label golang --label k8s --label architecture`,
		RunE: executeCreate,
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

	spinner := utils.StartSpinner("Creating tags…")
	createdTags, err := client.Tags.CreateOrUpdateBulk(cmd.Context(), inputs)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("create tags: %w", err)
	}

	renderTagsTable(createdTags, createdTags)

	return nil
}
