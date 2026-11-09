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

Without --collection the link goes to your Unorganized collection, like in the
web app. Without --name, Linkwarden uses the page's title.

To save many links at once from a JSON file, use "links import".`,
		Example: `  linkwarden links create https://go.dev/doc/effective_go
  linkwarden links create https://12factor.net --name "The Twelve-Factor App" --collection 3
  linkwarden links create https://k8spatterns.io --tag kubernetes --tag books`,
		Args: createArgs,
		RunE: executeCreate,
	}

	cmd.Flags().String("name", "", "Link name. Defaults to the page's title.")
	cmd.Flags().String("description", "", "Link description.")
	cmd.Flags().String("collection", "", "Collection ID to save to. Defaults to Unorganized.")
	cmd.Flags().StringSlice("tag", nil, "Tag name to add (repeatable).")

	// "links create --file" was the import command before "links import"
	// existed. Keep it working for scripts written against 1.0.
	cmd.Flags().String("file", "", "JSON file of links to import.")
	cmd.Flags().Int("concurrency", defaultConcurrency, "Number of links to create in parallel.")
	if err := cmd.Flags().MarkDeprecated("file", `use "links import --file" instead`); err != nil {
		return nil
	}
	if err := cmd.Flags().MarkHidden("concurrency"); err != nil {
		return nil
	}

	return cmd
}

// createArgs wants exactly one URL, or none when the deprecated --file is set.
func createArgs(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("file") {
		return cobra.NoArgs(cmd, args)
	}

	if len(args) != 1 {
		return fmt.Errorf("requires the URL to save, e.g. linkwarden links create https://example.com (got %d arguments)", len(args))
	}

	return nil
}

func executeCreate(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("file") {
		return executeImport(cmd, args)
	}

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

	collection, err := resolveCreateCollection(cmd.Flag("collection").Value.String(), collections)
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

	return utils.PrintRecord(link, link.ID)
}

// resolveCreateCollection picks the collection a single link is saved to: the
// one --collection names, else Unorganized, else whatever the picker returns.
// The collection must exist: an unknown ID would make the server create a
// fresh "Unorganized" collection for the link.
func resolveCreateCollection(flagValue string, collections []linkwarden.Collection) (linkwarden.Collection, error) {
	byID := make(map[int]linkwarden.Collection, len(collections))
	items := make([]utils.PickerItem, len(collections))
	for i, c := range collections {
		byID[c.ID] = c
		items[i] = utils.PickerItem{Label: c.Name, Value: c.ID}
	}

	if flagValue == "" {
		for _, c := range collections {
			if c.Name == unorganizedCollectionName && c.ParentID == nil {
				return c, nil
			}
		}
	}

	id, err := utils.ResolveID(
		flagValue,
		"--collection",
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
