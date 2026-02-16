package collections

import (
	"os"
	"path"
	"strconv"

	"github.com/CerealKiller97/linkwarden-cli/pkg/dto"
	"github.com/CerealKiller97/linkwarden-cli/pkg/services/collection"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	unorganizedCollectionName = "Unorganized"
	defaultColor              = "#00bba7"
	defaultDescription        = "Managed by: Linkwarden CLI"
)

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a collection",
		Long:  `Create a collection`,
		RunE:  executeCreate,
	}

	// Get unorganized collection ID from XDG_STATE_HOME/linkwarden/unorganized if exists
	getUnorganizedID := func() int {
		xdgStateHome, exists := os.LookupEnv("XDG_DATA_HOME")
		if !exists {
			xdgStateHome = os.Getenv("HOME") + "/.local/state"
		}

		// Read Unorganized collection ID from XDG_STATE_HOME/linkwarden/unorganized
		fetched, err := os.ReadFile(path.Join(xdgStateHome, "linkwarden", "unorganized"))
		if err != nil && !os.IsNotExist(err) {
			log.Err(err).Msg("Failed to read unorganized collection ID")
			return 0
		}

		unorganizedCollectionID := 0

		if len(fetched) > 0 {
			unorganizedCollectionID, err = strconv.Atoi(string(fetched))
			if err != nil {
				log.Err(err).Msg("Failed to convert unorganized collection ID to int")
				return 0
			}
		}

		return unorganizedCollectionID
	}

	cmd.Flags().String("name", "", "Collection name.")
	cmd.Flags().String("description", defaultDescription, "Collection description.")
	cmd.Flags().String("color", defaultColor, "Collection color.")
	cmd.Flags().String("icon", "", "Collection icon.")
	cmd.Flags().String("icon-weight", "", "Collection icon weight.")
	cmd.Flags().Int("parent-id", getUnorganizedID(), "Parent collection ID.")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}

	return cmd
}

func executeCreate(cmd *cobra.Command, _ []string) error {
	collectionService := collection.New(
		utils.GetInstanceURL(),
		utils.GetLinkwardenToken(),
		zerolog.New(os.Stdout),
	)

	name := cmd.Flag("name").Value.String()
	description := cmd.Flag("description").Value.String()
	color := cmd.Flag("color").Value.String()
	icon := cmd.Flag("icon").Value.String()
	iconWeight := cmd.Flag("icon-weight").Value.String()
	parentID := cmd.Flag("parent-id").Value.String()

	unorganizedCollectionID, err := strconv.Atoi(parentID)
	if err != nil {
		log.Err(err).Msg("Failed to convert parent ID to int")
		return err
	}

	// When unorganized collection ID is 0, find Unorganized collection
	if unorganizedCollectionID == 0 {
		collections, err := collectionService.Get(cmd.Context())
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

		// persist Unorganized collection ID to XDG_STATE_HOME/linkwarden/unorganized
		if unorganizedCollectionID != 0 {
			if err := persistUnorganizedCollectionID(unorganizedCollectionID); err != nil {
				return err
			}
		}
	}

	request := dto.CreateCollectionRequest{
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

	collection, err := collectionService.Create(cmd.Context(), request)
	if err != nil {
		log.Err(err).Msg("Failed to create collection")
		return err
	}

	log.Info().
		Int("id", collection.Response.ID).
		Msg("Collection created successfully")

	return nil
}

func persistUnorganizedCollectionID(unorganizedCollectionID int) error {
	xdgStateHome, exists := os.LookupEnv("XDG_STATE_HOME")

	// if not exist
	if !exists {
		xdgStateHome = os.Getenv("HOME") + "/.local/state"
	}

	if err := os.MkdirAll(path.Join(xdgStateHome, "/linkwarden"), 0755); err != nil {
		log.Err(err).Msg("Failed to create XDG_STATE_HOME/linkwarden/unorganized directory")
		return err
	}

	err := os.WriteFile(
		path.Join(xdgStateHome, "linkwarden", "unorganized"),
		[]byte(strconv.Itoa(unorganizedCollectionID)),
		0644,
	)
	if err != nil {
		return err
	}

	return nil
}
