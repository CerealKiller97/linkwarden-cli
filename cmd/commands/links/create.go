package links

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	errFileEmpty          = errors.New("file is empty")
	errCollectionNotFound = errors.New("collection not found")
)

const defaultConcurrency = 4

func Create() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "create",
		RunE: execute,
	}

	cmd.Flags().String("file", "", "File to upload.")
	cmd.Flags().String("collection", "", "Collection ID to upload to.")
	cmd.Flags().Int("concurrency", defaultConcurrency, "Number of links to create in parallel.")

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
	concurrency, err := cmd.Flags().GetInt("concurrency")
	if err != nil {
		return err
	}

	if concurrency < 1 {
		return fmt.Errorf("--concurrency must be at least 1, got %d", concurrency)
	}

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

	// An unknown ID would otherwise be sent as an empty collection, which makes
	// the server create a fresh "Unorganized" collection for every request.
	collection, ok := collectionByID[collectionID]
	if !ok {
		return fmt.Errorf("%w: %d", errCollectionNotFound, collectionID)
	}

	var (
		wg      sync.WaitGroup
		sem     = make(chan struct{}, concurrency)
		success atomic.Int64
		failed  atomic.Int64
	)

	bar := progressbar.Default(int64(total))

	for _, item := range data {
		sem <- struct{}{}
		wg.Go(func() {
			defer func() { <-sem }()

			_, err := client.Links.Create(cmd.Context(), linkwarden.CreateLinkRequest{
				Name:       item.Title,
				URL:        item.URL,
				Type:       "url",
				Collection: linkwarden.Collection{ID: collection.ID, Name: collection.Name},
				Tags:       []linkwarden.Tag{},
			})

			// A failed render of the progress bar shouldn't abort the import.
			_ = bar.Add(1)

			if err != nil {
				log.Err(err).
					Str("url", item.URL).
					Str("title", item.Title).
					Msg("Failed to create link")

				failed.Add(1)
				return
			}

			success.Add(1)
		})
	}

	wg.Wait()

	displayResults(total, int(success.Load()), int(failed.Load()))

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
