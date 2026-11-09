package links

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	errFileEmpty          = errors.New("file is empty")
	errCollectionNotFound = errors.New("collection not found")
	errLinksFailed        = errors.New("some links could not be created")
)

const defaultConcurrency = 4

func Import() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import links from a JSON file",
		Long: `Import links from a JSON file into one collection.

The file is a JSON array of objects with a "url" and a "title":

  [
    { "url": "https://go.dev/doc/effective_go", "title": "Effective Go" },
    { "url": "https://12factor.net", "title": "The Twelve-Factor App" }
  ]

Links already in the collection are skipped, as are URLs repeated within the
file. A summary of how many were created, skipped and failed is printed at the
end, and the command fails if any link could not be created.

Without --collection, a picker asks for the target collection.`,
		Example: `  linkwarden links import --file reading-list.json
  linkwarden links import --file reading-list.json --collection 3`,
		RunE: executeImport,
	}

	cmd.Flags().String("file", "", "JSON file of links to import.")
	cmd.Flags().String("collection", "", "Collection ID to import into.")
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

type failure struct {
	item item
	err  error
}

func executeImport(cmd *cobra.Command, _ []string) error {
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
		return fmt.Errorf("read file: %w", err)
	}

	var data []item
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return fmt.Errorf("unmarshal file: %w", err)
	}

	total := len(data)

	if total == 0 {
		return errFileEmpty
	}

	client := utils.NewLinkwardenClient()

	spinner := utils.StartSpinner("Fetching collections…")
	collections, err := client.Collections.List(cmd.Context())
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("get collections: %w", err)
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
		return fmt.Errorf("resolve collection: %w", err)
	}

	// An unknown ID would otherwise be sent as an empty collection, which makes
	// the server create a fresh "Unorganized" collection for every request.
	collection, ok := collectionByID[collectionID]
	if !ok {
		return fmt.Errorf("%w: %d", errCollectionNotFound, collectionID)
	}

	// Linkwarden happily stores the same URL twice unless the account has
	// "prevent duplicate links" on, so check the collection first.
	spinner = utils.StartSpinner("Checking for links already saved…")
	existing, err := existingURLs(cmd.Context(), client, collection.ID)
	spinner.Stop()
	if err != nil {
		return fmt.Errorf("list links already in the collection: %w", err)
	}

	toCreate, skipped := splitExisting(data, existing)

	var (
		wg            sync.WaitGroup
		sem           = make(chan struct{}, concurrency)
		success       atomic.Int64
		alreadyExists atomic.Int64
		mu            sync.Mutex
		failures      []failure
	)
	alreadyExists.Add(int64(skipped))

	// Nothing new to create: no point flashing an empty bar.
	bar := utils.NoProgress()
	if len(toCreate) > 0 {
		bar = utils.StartProgress("Creating links", len(toCreate))
	}

	ctx := cmd.Context()

	for _, item := range toCreate {
		sem <- struct{}{}
		// Ctrl+C: start no more requests; the ones in flight are cancelled too.
		if ctx.Err() != nil {
			<-sem
			break
		}
		wg.Go(func() {
			defer func() { <-sem }()

			_, err := client.Links.Create(ctx, linkwarden.CreateLinkRequest{
				Name:       item.Title,
				URL:        item.URL,
				Type:       "url",
				Collection: linkwarden.Collection{ID: collection.ID, Name: collection.Name},
				Tags:       []linkwarden.Tag{},
			})

			// Cut short by Ctrl+C: not a failure of this link, and logging one
			// line per cancelled request would bury the summary.
			if ctx.Err() != nil && errors.Is(err, context.Canceled) {
				return
			}

			bar.Increment()

			// Accounts with "prevent duplicate links" on reject them with a 409.
			if linkwarden.IsConflict(err) {
				alreadyExists.Add(1)
				return
			}

			if err != nil {
				mu.Lock()
				failures = append(failures, failure{item: item, err: err})
				mu.Unlock()

				return
			}

			success.Add(1)
		})
	}

	wg.Wait()
	bar.Stop()

	// Logged only once the bar is gone, so they don't tear through its redraws.
	for _, f := range failures {
		log.Err(f.err).
			Str("url", f.item.URL).
			Str("title", f.item.Title).
			Msg("Failed to create link")
	}

	displayResults(results{
		total:         total,
		success:       int(success.Load()),
		alreadyExists: int(alreadyExists.Load()),
		failed:        len(failures),
		failures:      failures,
	})

	// Interrupted: the summary above shows what got saved before Ctrl+C.
	if err := ctx.Err(); err != nil {
		return err
	}

	// Fail the command so scripts and cron jobs notice links that weren't saved.
	if len(failures) > 0 {
		return fmt.Errorf("%w: %d of %d", errLinksFailed, len(failures), total)
	}

	return nil
}

// splitExisting returns the items whose URL isn't in existing yet, and how many
// were skipped. A URL repeated within the file counts as already existing
// after its first occurrence. existing is updated in place.
func splitExisting(data []item, existing map[string]bool) ([]item, int) {
	toCreate := make([]item, 0, len(data))
	skipped := 0

	for _, it := range data {
		key := normalizeURL(it.URL)
		if existing[key] {
			skipped++
			continue
		}
		existing[key] = true
		toCreate = append(toCreate, it)
	}

	return toCreate, skipped
}

type results struct {
	total, success, alreadyExists, failed int
	failures                              []failure
}

// importReport is what --json prints for an import: the totals, and each link
// that failed with the reason.
type importReport struct {
	Total         int             `json:"total"`
	Success       int             `json:"success"`
	AlreadyExists int             `json:"alreadyExists"`
	Failed        int             `json:"failed"`
	Failures      []importFailure `json:"failures"`
}

type importFailure struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Error string `json:"error"`
}

func (r results) report() importReport {
	failures := make([]importFailure, len(r.failures))
	for i, f := range r.failures {
		failures[i] = importFailure{URL: f.item.URL, Title: f.item.Title, Error: f.err.Error()}
	}

	return importReport{
		Total:         r.total,
		Success:       r.success,
		AlreadyExists: r.alreadyExists,
		Failed:        r.failed,
		Failures:      failures,
	}
}

// Columns of the results table.
const (
	colTotal = iota
	colSuccess
	colAlreadyExists
	colFailed
)

func displayResults(r results) {
	t := utils.NewSummary("Import Results")
	t.Header("Total", "Success", "Already exists", "Failed")
	t.Row(r.total, r.success, r.alreadyExists, r.failed)
	t.JSON(r.report())
	t.Tone(func(_, col int, value string) utils.Tone {
		if value == "0" {
			return utils.ToneNone
		}

		switch col {
		case colSuccess:
			return utils.ToneSuccess
		case colAlreadyExists:
			return utils.ToneWarning
		case colFailed:
			return utils.ToneDanger
		default:
			return utils.ToneNone
		}
	})
	t.Render()
}
