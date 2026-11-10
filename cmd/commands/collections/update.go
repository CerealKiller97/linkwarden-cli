package collections

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Update() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a collection",
		RunE:  executeUpdate,
	}

	cmd.Flags().Int("id", 0, "Collection ID.")
	cmd.Flags().String("name", "", "New collection name.")
	cmd.Flags().String("description", "", "New collection description.")
	cmd.Flags().String("color", "", "New collection color.")
	cmd.Flags().String("icon", "", "New collection icon.")
	cmd.Flags().String("icon-weight", "", "New collection icon weight.")
	cmd.Flags().Int("parent-id", 0, "New parent collection ID.")
	cmd.Flags().Bool("public", false, "Make the collection public (--public=false makes it private).")

	return cmd
}

func executeUpdate(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd,
		"id",
		"Select Collection to update:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve collection ID: %w", err)
	}

	parentID, err := cmd.Flags().GetInt("parent-id")
	if err != nil {
		return err
	}
	var isPublic *bool
	if cmd.Flags().Changed("public") {
		v, err := cmd.Flags().GetBool("public")
		if err != nil {
			return err
		}
		isPublic = &v
	}

	spinner := utils.StartSpinner("Updating collection…")
	collection, err := client.Collections.Update(cmd.Context(), id, linkwarden.UpdateCollectionRequest{
		Name:        cmd.Flag("name").Value.String(),
		Description: cmd.Flag("description").Value.String(),
		Color:       cmd.Flag("color").Value.String(),
		Icon:        cmd.Flag("icon").Value.String(),
		IconWeight:  cmd.Flag("icon-weight").Value.String(),
		ParentID:    parentID,
		IsPublic:    isPublic,
	})
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("update collection: %w", err)
	}

	log.Info().Int("id", collection.ID).Msg("Collection updated successfully")

	return utils.PrintRecord(collection, collection.ID)
}
