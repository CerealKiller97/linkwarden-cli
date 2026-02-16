package links

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/CerealKiller97/linkwarden-cli/pkg/dto"
	"github.com/CerealKiller97/linkwarden-cli/pkg/services/collection"
	"github.com/CerealKiller97/linkwarden-cli/pkg/services/link"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var errFileEmpty = errors.New("file is empty")

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		RunE: execute,
	}

	// TODO: Help
	// cmd.SetHelpFunc(func(c *cobra.Command, e []string) {})
	// cmd.SetHelpTemplate()
	// cmd.SetHelpTemplate()

	cmd.Flags().String("file", "", "File to upload.")
	cmd.Flags().String("collection", "", "Collection to upload to.")

	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil
	}

	return cmd
}

type item struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

func execute(cmd *cobra.Command, _ []string) error {
	filePath := cmd.Flag("file").Value.String()
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Err(err).Msg("failed to read file")
		return err
	}

	var data []item
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		log.Err(err).Msg("failed to unmarshal file")
		return err
	}

	total := len(data)

	if total == 0 {
		return errFileEmpty
	}

	collectionService := collection.New(
		utils.GetInstanceURL(),
		utils.GetLinkwardenToken(),
		zerolog.New(os.Stdout),
	)

	collections, err := collectionService.Get(cmd.Context())
	if err != nil {
		log.Err(err).Msg("Failed to get collections")
		return err
	}

	collectionMap := make(map[string]int, len(collections))
	collectionIDTONameMap := make(map[int]string, len(collections))
	collectionNames := make([]string, 0, len(collections))

	for _, item := range collections {
		collectionMap[item.Name] = item.ID
		collectionIDTONameMap[item.ID] = item.Name
		collectionNames = append(collectionNames, item.Name)
	}

	collectionID := cmd.Flag("collection").Value.String()
	collection := ""

	// Show interactive prompt if a collection is not specified
	if collectionID == "" {
		prompt := promptui.Select{
			Label: "Select Collection where you want to upload links to:",
			Items: collectionNames,
		}

		_, collection, err = prompt.Run()
		if err != nil {
			log.Err(err).Msg("Prompt Failed")
			return err
		}
	}

	success := 0
	failed := 0

	linkService := link.New(
		utils.GetInstanceURL(),
		utils.GetLinkwardenToken(),
		zerolog.New(os.Stdout),
	)

	bar := progressbar.Default(100)

	stepSize := 100.0 / float64(len(data))

	for _, item := range data {
		var coll dto.Collection
		if collection == "" {
			collectionIDInt, err := strconv.Atoi(collectionID)
			if err != nil {
				return err
			}
			coll = dto.Collection{
				ID:   collectionIDInt,
				Name: collectionIDTONameMap[collectionIDInt],
			}
		} else {
			coll = dto.Collection{
				ID:   collectionMap[collection],
				Name: collection,
			}
		}

		_, err = linkService.Create(cmd.Context(), dto.CreateLinkRequest{
			Name:        item.Title,
			URL:         item.URL,
			Type:        "url",
			Description: "",
			Collection:  coll,
			Tags:        []dto.Tag{},
		})
		if err != nil {
			log.Err(err).
				Str("url", item.URL).
				Str("title", item.Title).
				Msg("Failed to create link")

			failed++
			continue
		}

		if err := bar.Add(int(stepSize)); err != nil {
			return err
		}

		success++
	}

	displayResults(total, success, failed)

	return nil
}

func displayResults(total int, success int, failed int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Linkwarden")
	t.AppendHeader(table.Row{"Total", "Success", "Failed"})
	t.AppendRows([]table.Row{
		{total, success, failed},
	})
	t.SetStyle(table.StyleColoredBlackOnCyanWhite)
	t.AppendSeparator()
	t.Render()
}
