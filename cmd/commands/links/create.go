package links

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/jedib0t/go-pretty/v6/table"
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

	cmd.Flags().String("file", "", "File to upload.")
	cmd.Flags().String("collection", "", "Collection ID to upload to.")

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

	client := utils.NewLinkwardenClient()

	collections, err := client.Collections.List(cmd.Context())
	if err != nil {
		log.Err(err).Msg("Failed to get collections")
		return err
	}

	collectionByID := make(map[int]linkwarden.Collection, len(collections))
	items := make([]utils.PickerItem, len(collections))
	for i, c := range collections {
		collectionByID[c.ID] = c
		items[i] = utils.PickerItem{Label: c.Name, Value: c.ID}
	}

	collectionID, err := utils.ResolveID(
		cmd.Flag("collection").Value.String(),
		"--collection",
		"Select Collection where you want to upload links to:",
		func() ([]utils.PickerItem, error) { return items, nil },
	)
	if err != nil {
		log.Err(err).Msg("Failed to resolve collection")
		return err
	}

	collection := collectionByID[collectionID]

	success := 0
	failed := 0

	bar := progressbar.Default(100)

	stepSize := 100.0 / float64(len(data))

	for _, item := range data {
		_, err = client.Links.Create(cmd.Context(), linkwarden.CreateLinkRequest{
			Name:       item.Title,
			URL:        item.URL,
			Type:       "url",
			Collection: linkwarden.Collection{ID: collection.ID, Name: collection.Name},
			Tags:       []linkwarden.Tag{},
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
	t := utils.NewTable("Linkwarden")
	t.AppendHeader(table.Row{"Total", "Success", "Failed"})
	t.AppendRows([]table.Row{
		{total, success, failed},
	})
	t.AppendSeparator()
	t.Render()
}
