package collections

import (
	"log/slog"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func Update() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a collection",
		RunE:  executeUpdate,
	}

	cmd.Flags().String("id", "", "Collection ID.")
	cmd.Flags().String("name", "", "New collection name.")
	cmd.Flags().String("description", "", "New collection description.")
	cmd.Flags().String("color", "", "New collection color.")
	cmd.Flags().String("icon", "", "New collection icon.")
	cmd.Flags().String("icon-weight", "", "New collection icon weight.")
	cmd.Flags().Int("parent-id", 0, "New parent collection ID.")
	cmd.Flags().Bool("public", false, "Make the collection public.")

	return cmd
}

func executeUpdate(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Collection to update:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		slog.Error("Failed to resolve collection ID", slog.Any("error", err))
		return err
	}

	parentID, err := cmd.Flags().GetInt("parent-id")
	if err != nil {
		return err
	}
	isPublic, err := cmd.Flags().GetBool("public")
	if err != nil {
		return err
	}

	collection, err := client.Collections.Update(cmd.Context(), id, linkwarden.UpdateCollectionRequest{
		Name:        cmd.Flag("name").Value.String(),
		Description: cmd.Flag("description").Value.String(),
		Color:       cmd.Flag("color").Value.String(),
		Icon:        cmd.Flag("icon").Value.String(),
		IconWeight:  cmd.Flag("icon-weight").Value.String(),
		ParentID:    parentID,
		IsPublic:    isPublic,
	})
	if err != nil {
		slog.Error("Failed to update collection", slog.Any("error", err))
		return err
	}

	slog.Info("Collection updated successfully", slog.Int("id", collection.ID))

	return nil
}
