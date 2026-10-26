package collections

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	unorganizedCollectionName = "Unorganized"
	defaultColor              = "#00bba7"
	defaultDescription        = "Managed by: Linkwarden CLI"
	unorganizedStateFile      = "unorganized"
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

	unorganizedCollectionID, err := cmd.Flags().GetInt("parent-id")
	if err != nil {
		return err
	}

	if !cmd.Flags().Changed("parent-id") {
		unorganizedCollectionID = readUnorganizedCollectionID()
	}

	// When unorganized collection ID is 0, find Unorganized collection
	if unorganizedCollectionID == 0 {
		collections, err := client.Collections.List(cmd.Context())
		if err != nil {
			log.Err(err).Msg("Failed to get collections")
			return err
		}

		for _, collection := range collections {
			if collection.Name == name {
				log.Err(err).Msg("Collection already exists")
				return err
			}

			// find Unorganized collection
			if collection.Name == unorganizedCollectionName {
				unorganizedCollectionID = collection.ID
				break
			}
		}

		// cache the Unorganized collection ID in $XDG_STATE_HOME/linkwarden/unorganized
		if unorganizedCollectionID != 0 {
			if err := persistUnorganizedCollectionID(unorganizedCollectionID); err != nil {
				return err
			}
		}
	}

	request := linkwarden.CreateCollectionRequest{
		Name:        name,
		Description: description,
		Color:       color,
		Icon:        icon,
		IconWeight:  iconWeight,
		ParentID:    unorganizedCollectionID,
	}

	if request.ParentID == 0 && unorganizedCollectionID != 0 {
		request.ParentID = unorganizedCollectionID
	}

	collection, err := client.Collections.Create(cmd.Context(), request)
	if err != nil {
		log.Err(err).Msg("Failed to create collection")
		return err
	}

	log.Info().
		Int("id", collection.ID).
		Msg("Collection created successfully")

	return nil
}

// readUnorganizedCollectionID returns the Unorganized collection ID cached in
// $XDG_STATE_HOME/linkwarden/unorganized, or 0 if there is no usable cached value.
func readUnorganizedCollectionID() int {
	file, err := utils.StateFile(unorganizedStateFile)
	if err != nil {
		log.Err(err).Msg("Failed to resolve state directory")
		return 0
	}

	fetched, err := os.ReadFile(file)
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

func persistUnorganizedCollectionID(unorganizedCollectionID int) error {
	file, err := utils.StateFile(unorganizedStateFile)
	if err != nil {
		log.Err(err).Msg("Failed to resolve state directory")
		return err
	}

	if err := os.MkdirAll(filepath.Dir(file), 0o750); err != nil {
		log.Err(err).Msg("Failed to create state directory")
		return err
	}

	return os.WriteFile(file, []byte(strconv.Itoa(unorganizedCollectionID)), 0o600)
}
