package links

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
		return fmt.Errorf("resolve link ID: %w", err)
	}

	link, err := client.Links.Update(cmd.Context(), id, linkwarden.UpdateLinkRequest{
		Name:        cmd.Flag("name").Value.String(),
		URL:         cmd.Flag("url").Value.String(),
		Description: cmd.Flag("description").Value.String(),
		Icon:        cmd.Flag("icon").Value.String(),
		IconWeight:  cmd.Flag("icon-weight").Value.String(),
		Color:       cmd.Flag("color").Value.String(),
	})
	if err != nil {
		return fmt.Errorf("update link: %w", err)
	}

	log.Info().Int("id", link.ID).Msg("Link updated successfully")

	return nil
}
