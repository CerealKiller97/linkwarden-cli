package collections

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Msg("Failed to resolve collection ID")
		return err
	}

	// Start from the current collection: the API requires name and members and
	// replaces members with whatever is sent.
	current, err := client.Collections.Get(cmd.Context(), id)
	if err != nil {
		log.Err(err).Msg("Failed to get collection")
		return err
	}

	req := linkwarden.UpdateCollectionRequest{
		Name:        current.Name,
		Description: current.Description,
		Color:       current.Color,
		Icon:        current.Icon,
		IconWeight:  current.IconWeight,
		ParentID:    current.ParentID,
		IsPublic:    &current.IsPublic,
		Members:     current.Members,
	}

	flags := cmd.Flags()
	for flag, field := range map[string]*string{
		"name":        &req.Name,
		"description": &req.Description,
		"color":       &req.Color,
		"icon":        &req.Icon,
		"icon-weight": &req.IconWeight,
	} {
		if flags.Changed(flag) {
			*field = cmd.Flag(flag).Value.String()
		}
	}

	if flags.Changed("parent-id") {
		parentID, err := flags.GetInt("parent-id")
		if err != nil {
			return err
		}
		req.ParentID = &parentID
	}

	if flags.Changed("public") {
		isPublic, err := flags.GetBool("public")
		if err != nil {
			return err
		}
		req.IsPublic = &isPublic
	}

	collection, err := client.Collections.Update(cmd.Context(), id, req)
	if err != nil {
		log.Err(err).Msg("Failed to update collection")
		return err
	}

	log.Info().Int("id", collection.ID).Msg("Collection updated successfully")

	return nil
}
