package collections

import (
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Delete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a collection",
		RunE:  executeDelete,
	}

	cmd.Flags().String("id", "", "Collection ID.")
	utils.AddYesFlag(cmd)

	return cmd
}

func executeDelete(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	id, err := utils.ResolveID(
		cmd.Flag("id").Value.String(),
		"--id",
		"Select Collection to delete:",
		collectionsPickerFetch(cmd.Context(), client),
	)
	if err != nil {
		return fmt.Errorf("resolve collection ID: %w", err)
	}

	// Deleting a collection also deletes every link in it, so name both. The
	// lookup only happens when someone is actually asked.
	ok, err := utils.ConfirmWith(cmd, func() (string, error) {
		spinner := utils.StartSpinner("Fetching collection…")
		collection, err := client.Collections.Get(cmd.Context(), id)
		spinner.Stop()
		if err != nil {
			return "", fmt.Errorf("get collection: %w", err)
		}

		return deleteQuestion(collection), nil
	})
	if err != nil {
		return err
	}
	if !ok {
		log.Info().Msg("Cancelled; nothing was deleted")
		return nil
	}

	spinner := utils.StartSpinner("Deleting collection…")
	err = client.Collections.Delete(cmd.Context(), id)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}

	log.Info().Int("id", id).Msg("Collection deleted successfully")

	return nil
}

// deleteQuestion asks to delete c, naming how many links go with it when the
// API reported a count.
func deleteQuestion(c *linkwarden.Collection) string {
	if c.Count == nil {
		return fmt.Sprintf("Delete collection %q and all its links? This can't be undone.", c.Name)
	}

	return fmt.Sprintf("Delete collection %q and its %s? This can't be undone.",
		c.Name, utils.Plural(c.Count.Links, "link"))
}
