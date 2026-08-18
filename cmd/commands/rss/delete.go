package rss

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an RSS subscription",
		RunE:  executeDelete,
	}

	cmd.Flags().String("id", "", "Subscription ID.")

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select RSS subscription to delete:",
		rssPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve RSS subscription ID")
		return err
	}

	if err := client.RSS.Delete(cmd.Context(), id); err != nil {
		log.Err(err).Msg("Failed to delete RSS subscription")
		return err
	}

	log.Info().Int("id", id).Msg("RSS subscription deleted successfully")

	return nil
}
