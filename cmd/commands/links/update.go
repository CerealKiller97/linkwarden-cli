package links

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Update() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a link",
		RunE:  executeUpdate,
	}

	cmd.Flags().String("id", "", "Link ID.")
	cmd.Flags().String("name", "", "New link name.")
	cmd.Flags().String("url", "", "New link URL.")
	cmd.Flags().String("description", "", "New link description.")
	cmd.Flags().String("icon", "", "New link icon.")
	cmd.Flags().String("icon-weight", "", "New link icon weight.")
	cmd.Flags().String("color", "", "New link color.")

	return cmd
}

func executeUpdate(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Link to update:",
		linksPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve link ID")
		return err
	}

	// Start from the current link: the API replaces the whole link, so any
	// field left out of the request would be cleared.
	current, err := client.Links.Get(cmd.Context(), id)
	if err != nil {
		log.Err(err).Msg("Failed to get link")
		return err
	}

	req := linkwarden.UpdateLinkRequest{
		Name:        current.Name,
		URL:         current.URL,
		Description: current.Description,
		Icon:        current.Icon,
		IconWeight:  current.IconWeight,
		Color:       current.Color,
		Tags:        current.Tags,
		Collection:  linkwarden.Collection{ID: current.Collection.ID, OwnerID: current.Collection.OwnerID},
		PinnedBy:    current.PinnedBy,
	}

	for flag, field := range map[string]*string{
		"name":        &req.Name,
		"url":         &req.URL,
		"description": &req.Description,
		"icon":        &req.Icon,
		"icon-weight": &req.IconWeight,
		"color":       &req.Color,
	} {
		if cmd.Flags().Changed(flag) {
			*field = cmd.Flag(flag).Value.String()
		}
	}

	link, err := client.Links.Update(cmd.Context(), id, req)
	if err != nil {
		log.Err(err).Msg("Failed to update link")
		return err
	}

	log.Info().Int("id", link.ID).Msg("Link updated successfully")

	return nil
}
