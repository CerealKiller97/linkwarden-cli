package pickers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestShortURL(t *testing.T) {
	assert := require.New(t)

	assert.Equal("go.dev/doc/effective_go", ShortURL("https://go.dev/doc/effective_go"))
	assert.Equal("example.com", ShortURL("https://www.example.com/"))
	assert.Equal("example.com/search?q=go", ShortURL("http://example.com/search?q=go"))
	assert.Equal("stolicefotelje.rs/namestaj/radne-stolic…", ShortURL("https://stolicefotelje.rs/namestaj/radne-stolice/radna-stolica-mirra/"))
	assert.Equal("not a url", ShortURL("not a url"))
	assert.Empty(ShortURL(""))
}

// fakeLinks serves ids newest first, two per page, by cursor like the API. It
// answers both the old /api/v1/links and the /api/v1/search route, with both
// response shapes, so it works whichever the SDK calls.
func fakeLinks(t *testing.T, ids ...int) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cursor, _ := strconv.Atoi(r.URL.Query().Get("cursor"))
		page := ""
		n := 0
		for _, id := range ids {
			if (cursor != 0 && id >= cursor) || n == 2 {
				continue
			}
			if page != "" {
				page += ","
			}
			page += `{"id":` + strconv.Itoa(id) + `,"name":"","url":"https://example.com/` + strconv.Itoa(id) + `"}`
			n++
		}
		_, _ = w.Write([]byte(`{"response":[` + page + `],"data":{"links":[` + page + `]}}`))
	}))
}

func TestLinksPagesByLastID(t *testing.T) {
	assert := require.New(t)

	server := fakeLinks(t, 5, 4, 3)
	defer server.Close()
	fetch := Links(t.Context(), linkwarden.New(server.URL, "test-token"))

	items, next, err := fetch(0)
	assert.NoError(err)
	assert.Equal(4, next, "the next cursor is the last ID on the page")
	assert.Equal(5, items[0].Value)
	assert.Equal("https://example.com/5", items[0].Label, "a link without a name shows its URL")
	assert.Equal("example.com/5", items[0].Detail)

	items, next, err = fetch(next)
	assert.NoError(err)
	assert.Len(items, 1)
	assert.Equal(3, next)

	items, next, err = fetch(next)
	assert.NoError(err)
	assert.Empty(items)
	assert.Zero(next, "an empty page ends paging")
}
