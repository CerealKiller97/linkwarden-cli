package links

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUploadsLinksFromFile(t *testing.T) {
	assert := require.New(t)

	var createCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/search":
			// With no matches the search route sends "data": [] instead of {"links": [...]}.
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"Success","data":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/links":
			createCount++
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":1,"name":"x","url":"https://example.com"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	file := filepath.Join(t.TempDir(), "links.json")
	assert.NoError(os.WriteFile(file, []byte(`[{"url":"https://example.com","title":"Example"}]`), 0o644))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})
	assert.NoError(cmd.Execute())
	assert.Equal(1, createCount)
}

func TestCreateEmptyFileErrors(t *testing.T) {
	assert := require.New(t)

	file := filepath.Join(t.TempDir(), "links.json")
	assert.NoError(os.WriteFile(file, []byte(`[]`), 0o644))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})
	err := cmd.Execute()
	assert.ErrorIs(err, errFileEmpty)
}

func TestCreateUnknownCollectionErrorsWithoutCreatingLinks(t *testing.T) {
	assert := require.New(t)

	var createCount atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"}]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/links":
			createCount.Add(1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"response":{"id":1,"name":"x","url":"https://example.com"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	file := filepath.Join(t.TempDir(), "links.json")
	assert.NoError(os.WriteFile(file, []byte(`[{"url":"https://example.com","title":"Example"}]`), 0o644))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "9999"})
	err := cmd.Execute()
	assert.ErrorIs(err, errCollectionNotFound)
	assert.ErrorContains(err, "9999")
	assert.Zero(createCount.Load())
}

// fakeLinkServer serves one collection whose links are pre-seeded with
// existing URLs, paged pageSize at a time by cursor like the real API, and
// counts creates. Creating a URL in conflict returns 409.
func fakeLinkServer(t *testing.T, existing []string, pageSize int, conflict string) (*httptest.Server, *atomic.Int64, *atomic.Int64) {
	t.Helper()

	var creates, lists atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Go"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/search":
			lists.Add(1)
			assert.Equal(t, "1", r.URL.Query().Get("collectionId"))
			// IDs count down from len(existing), newest first, like the API.
			cursor, _ := strconv.Atoi(r.URL.Query().Get("cursor"))
			var page []string
			for i := len(existing); i >= 1 && len(page) < pageSize; i-- {
				if cursor != 0 && i >= cursor {
					continue
				}
				page = append(page, fmt.Sprintf(`{"id":%d,"url":%q}`, i, existing[i-1]))
			}
			if len(page) == 0 {
				_, _ = w.Write([]byte(`{"message":"Success","data":[]}`))
				return
			}
			_, _ = w.Write([]byte(`{"message":"Success","data":{"links":[` + strings.Join(page, ",") + `]}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/links":
			body, _ := io.ReadAll(r.Body)
			if conflict != "" && strings.Contains(string(body), conflict) {
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(`{"response":"Link already exists"}`))
				return
			}
			creates.Add(1)
			_, _ = w.Write([]byte(`{"response":{"id":100,"url":"x"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	utils.SetVars("test-token", server.URL)

	return server, &creates, &lists
}

func writeLinksFile(t *testing.T, urls ...string) string {
	t.Helper()

	items := make([]string, len(urls))
	for i, u := range urls {
		items[i] = fmt.Sprintf(`{"url":%q,"title":"t"}`, u)
	}
	file := filepath.Join(t.TempDir(), "links.json")
	require.NoError(t, os.WriteFile(file, []byte("["+strings.Join(items, ",")+"]"), 0o644))

	return file
}

func TestCreateSkipsLinksAlreadyInCollection(t *testing.T) {
	assert := require.New(t)

	// Five existing links served two per page, so the check has to paginate.
	existing := []string{"https://a.example", "https://b.example", "https://c.example", "https://d.example", "https://e.example/"}
	_, creates, lists := fakeLinkServer(t, existing, 2, "")

	file := writeLinksFile(t,
		"https://A.example/",  // exists (case and trailing slash differ)
		"https://e.example",   // exists (stored with a trailing slash)
		"https://new.example", // new
		"https://new.example", // repeated in the file
	)

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})
	assert.NoError(cmd.Execute())

	assert.Equal(int64(1), creates.Load(), "only the one new URL is created")
	assert.Equal(int64(4), lists.Load(), "3 pages of links plus the empty page that ends paging")
}

func TestCreateConflictCountsAsAlreadyExists(t *testing.T) {
	assert := require.New(t)

	_, creates, _ := fakeLinkServer(t, nil, 50, "dup.example")

	file := writeLinksFile(t, "https://dup.example", "https://ok.example")
	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})
	assert.NoError(cmd.Execute())
	assert.Equal(int64(1), creates.Load())
}

func TestCreateFailedLinksReturnError(t *testing.T) {
	assert := require.New(t)

	var creates atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Go"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/search":
			// No links in the collection yet: the search route sends "data": [].
			_, _ = w.Write([]byte(`{"message":"Success","data":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/links":
			body, _ := io.ReadAll(r.Body)
			if strings.Contains(string(body), "broken.example") {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"response":"boom"}`))
				return
			}
			creates.Add(1)
			_, _ = w.Write([]byte(`{"response":{"id":100,"url":"x"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	utils.SetVars("test-token", server.URL)

	file := writeLinksFile(t, "https://ok.example", "https://broken.example")
	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})

	err := cmd.Execute()
	assert.ErrorIs(err, errLinksFailed)
	assert.ErrorContains(err, "1 of 2")
	assert.Equal(int64(1), creates.Load(), "the working link is still created")
}

func TestSplitExisting(t *testing.T) {
	assert := require.New(t)

	existing := map[string]bool{"https://a.example": true}
	toCreate, skipped := splitExisting([]item{
		{URL: "https://a.example/"},
		{URL: "https://b.example"},
		{URL: " https://B.example/ "},
	}, existing)

	assert.Equal(2, skipped)
	assert.Equal([]item{{URL: "https://b.example"}}, toCreate)
}

func TestNormalizeURL(t *testing.T) {
	assert := require.New(t)

	assert.Equal("https://example.com", normalizeURL(" HTTPS://Example.COM/ "))
	assert.Equal("https://example.com/Docs", normalizeURL("https://example.com/Docs/"), "path case is kept")
	assert.Equal("https://example.com/a?x=1#y", normalizeURL("https://example.com/a?x=1#y"))
	assert.Equal("not a url", normalizeURL("not a url/"))
}

func TestCreateStopsWhenInterrupted(t *testing.T) {
	assert := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Unblocks hanging handlers if the test fails before the client hangs up.
	stopServer := make(chan struct{})

	var creates atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Unorganized"}]}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/search":
			_, _ = w.Write([]byte(`{"message":"Success","data":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/links":
			// The first request "presses Ctrl+C" and then hangs until cancelled.
			creates.Add(1)
			cancel()
			select {
			case <-r.Context().Done():
			case <-stopServer:
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()
	defer close(stopServer)
	utils.SetVars("test-token", server.URL)

	var items []string
	for i := range 50 {
		items = append(items, fmt.Sprintf(`{"url":"https://example.com/%d","title":"%d"}`, i, i))
	}
	file := filepath.Join(t.TempDir(), "links.json")
	assert.NoError(os.WriteFile(file, []byte("["+strings.Join(items, ",")+"]"), 0o644))

	cmd := Create()
	cmd.SetArgs([]string{"--file", file, "--collection", "1"})
	err := cmd.ExecuteContext(ctx)

	assert.ErrorIs(err, context.Canceled)
	assert.NotErrorIs(err, errLinksFailed, "cancelled links aren't failures")
	assert.LessOrEqual(creates.Load(), int64(defaultConcurrency), "no new requests after Ctrl+C")
}
