package collections

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	unorganizedCollectionName = "Unorganized"
	defaultColor              = "#00bba7"
	defaultDescription        = "Managed by: Linkwarden CLI"
	unorganizedStateDir       = "unorganized"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a collection",
		Long:  `Create a collection`,
		RunE:  executeCreate,
	}

	cmd.Flags().String("name", "", "Collection name.")
	cmd.Flags().String("description", defaultDescription, "Collection description.")
	cmd.Flags().String("color", defaultColor, "Collection color.")
	cmd.Flags().String("icon", "", "Collection icon.")
	cmd.Flags().String("icon-weight", "", "Collection icon weight.")
	cmd.Flags().Int("parent-id", 0, "Parent collection ID. Defaults to the Unorganized collection.")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}

	return cmd
}

func executeCreate(cmd *cobra.Command, _ []string) error {
	client := utils.NewLinkwardenClient()

	name := cmd.Flag("name").Value.String()
	description := cmd.Flag("description").Value.String()
	color := cmd.Flag("color").Value.String()
	icon := cmd.Flag("icon").Value.String()
	iconWeight := cmd.Flag("icon-weight").Value.String()

	parentID, err := cmd.Flags().GetInt("parent-id")
	if err != nil {
		return err
	}

	// Without --parent-id the collection goes under the account's Unorganized
	// collection. Its ID is cached per instance, since it differs between them.
	if !cmd.Flags().Changed("parent-id") {
		instanceURL := utils.GetInstanceURL()

		parentID = readUnorganizedCollectionID(instanceURL)
		if parentID == 0 {
			spinner := utils.StartSpinner("Checking existing collections…")
			collections, err := client.Collections.List(cmd.Context())
			spinner.Stop()
			if err != nil {
				return fmt.Errorf("get collections: %w", err)
			}

			parentID, err = unorganizedParentID(collections, name)
			if err != nil {
				return err
			}

			if parentID != 0 {
				if err := persistUnorganizedCollectionID(instanceURL, parentID); err != nil {
					return err
				}
			}
		}
	}

	request := linkwarden.CreateCollectionRequest{
		Name:        name,
		Description: description,
		Color:       color,
		Icon:        icon,
		IconWeight:  iconWeight,
		ParentID:    parentID,
	}

	spinner := utils.StartSpinner("Creating collection…")
	collection, err := client.Collections.Create(cmd.Context(), request)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	log.Info().
		Int("id", collection.ID).
		Msg("Collection created successfully")

	return utils.PrintRecord(collection, collection.ID)
}

// unorganizedParentID returns the ID of the Unorganized collection in
// collections, or 0 (a top-level collection) if there is none. It fails if a
// collection called name already exists.
func unorganizedParentID(collections []linkwarden.Collection, name string) (int, error) {
	for _, collection := range collections {
		if collection.Name == name {
			return 0, fmt.Errorf("collection %q already exists", name)
		}
	}

	return unorganizedID(collections), nil
}

// readUnorganizedCollectionID returns the Unorganized collection ID cached for
// instanceURL in $XDG_STATE_HOME/linkwarden/unorganized/<instance>, or 0 if
// there is no usable cached value.
func readUnorganizedCollectionID(instanceURL string) int {
	file, err := utils.InstanceStateFile(unorganizedStateDir, instanceURL)
	if err != nil {
		log.Err(err).Msg("Failed to resolve state directory")
		return 0
	}

	// Older versions kept a single ID for every instance in a file where the
	// per-instance directory now goes. It can't be trusted, so it's a miss.
	if info, err := os.Stat(filepath.Dir(file)); err != nil || !info.IsDir() {
		return 0
	}

	fetched, err := os.ReadFile(file) //nolint:gosec // path is the CLI's own state file
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Err(err).Msg("Failed to read unorganized collection ID")
		}
		return 0
	}

	id, err := strconv.Atoi(strings.TrimSpace(string(fetched)))
	if err != nil {
		log.Err(err).Msg("Failed to convert unorganized collection ID to int")
		return 0
	}

	return id
}

// unorganizedID returns the ID of the Unorganized collection in collections, or
// 0 if there is none.
func unorganizedID(collections []linkwarden.Collection) int {
	for _, collection := range collections {
		if collection.Name == unorganizedCollectionName {
			return collection.ID
		}
	}

	return 0
}

// persistUnorganizedCollectionID caches id as the Unorganized collection ID for
// instanceURL, replacing the single all-instances file older versions wrote.
func persistUnorganizedCollectionID(instanceURL string, id int) error {
	file, err := utils.InstanceStateFile(unorganizedStateDir, instanceURL)
	if err != nil {
		return fmt.Errorf("resolve state directory: %w", err)
	}

	dir := filepath.Dir(file)
	if info, err := os.Stat(dir); err == nil && !info.IsDir() {
		if err := os.Remove(dir); err != nil {
			return fmt.Errorf("remove old unorganized collection cache: %w", err)
		}
	}

	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create state directory: %w", err)
	}

	return os.WriteFile(file, []byte(strconv.Itoa(id)), 0o600)
}
