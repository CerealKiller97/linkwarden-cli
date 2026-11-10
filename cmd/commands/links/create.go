package links

import (
	"errors"
	"fmt"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var errLinkExists = errors.New("link already exists")

// unorganizedCollectionName is the collection Linkwarden saves links to when
// none is chosen.
const unorganizedCollectionName = "Unorganized"

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <url>",
		Short: "Save a link",
		Long: `Save a single link.

Without --collection-id the link goes to your Unorganized collection, like in the
web app. Without --name, Linkwarden uses the page's title.

To save many links at once from a JSON file, use "links import".`,
		Example: `  # Save a link to Unorganized; its name comes from the page's title
  linkwarden links create https://go.dev/doc/effective_go

  # Name it and save it to collection 3
  linkwarden links create https://12factor.net --name "The Twelve-Factor App" --collection-id 3

  # Add tags
  linkwarden links create https://k8spatterns.io --tag kubernetes --tag books`,
		Args: createArgs,
		RunE: executeCreate,
	}

	cmd.Flags().String("name", "", "Link name. Defaults to the page's title.")
	cmd.Flags().String("description", "", "Link description.")
	cmd.Flags().Int("collection-id", 0, "Collection ID to save to. Defaults to Unorganized.")
	cmd.Flags().StringSlice("tag", nil, "Tag name to add (repeatable).")

	return cmd
}

// createArgs wants exactly one URL, and points file imports at links import.
func createArgs(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(`requires the URL to save, e.g. linkwarden links create https://example.com (got %d arguments); to import a JSON file use "links import --file"`, len(args))
	}

	return nil
}

func executeCreate(cmd *cobra.Command, args []string) error {
	tagNames, err := cmd.Flags().GetStringSlice("tag")
	if err != nil {
		return err
	}

	tags := make([]linkwarden.Tag, len(tagNames))
	for i, name := range tagNames {
		tags[i] = linkwarden.Tag{Name: name}
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching collections…")
	collections, err := client.Collections.List(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get collections: %w", err)
	}

	collection, err := resolveCreateCollection(cmd, collections)
	if err != nil {
		return err
	}

	spinner = utils.StartSpinner("Saving link…")
	link, err := client.Links.Create(cmd.Context(), linkwarden.CreateLinkRequest{
		Name:        cmd.Flag("name").Value.String(),
		URL:         args[0],
		Type:        "url",
		Description: cmd.Flag("description").Value.String(),
		Tags:        tags,
		Collection:  linkwarden.Collection{ID: collection.ID, Name: collection.Name},
	})
	spinner.Stop()

	// Accounts with "prevent duplicate links" on reject them with a 409.
	if linkwarden.IsConflict(err) {
		return fmt.Errorf("%w: %s", errLinkExists, args[0])
	}
	if err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	log.Info().Int("id", link.ID).Msg("Link created successfully")

	return utils.PrintRecord(cmd.Context(), link, link.ID)
}

// resolveCreateCollection picks the collection a single link is saved to: the
// one --collection-id names, else Unorganized, else whatever the picker returns.
// The collection must exist: an unknown ID would make the server create a
// fresh "Unorganized" collection for the link.
func resolveCreateCollection(cmd *cobra.Command, collections []linkwarden.Collection) (linkwarden.Collection, error) {
	byID := make(map[int]linkwarden.Collection, len(collections))
	items := make([]utils.PickerItem, len(collections))
	for i, c := range collections {
		byID[c.ID] = c
		items[i] = utils.PickerItem{Label: c.Name, Value: c.ID}
	}

	if !cmd.Flags().Changed("collection-id") {
		for _, c := range collections {
			if c.Name == unorganizedCollectionName && c.ParentID == nil {
				return c, nil
			}
		}
	}

	id, err := utils.ResolveID(
		cmd,
		"collection-id",
		"Select Collection to save the link to:",
		func() ([]utils.PickerItem, error) { return items, nil },
	)
	if err != nil {
		return linkwarden.Collection{}, fmt.Errorf("resolve collection: %w", err)
	}

	collection, ok := byID[id]
	if !ok {
		return linkwarden.Collection{}, fmt.Errorf("%w: %d", errCollectionNotFound, id)
	}

	return collection, nil
}
