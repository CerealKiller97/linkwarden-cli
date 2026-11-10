package links

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

type createdLink struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Tags        []struct {
		Name string `json:"name"`
	} `json:"tags"`
	Collection struct {
		ID int `json:"id"`
	} `json:"collection"`
}

// fakeCreateServer serves collections and records each link POSTed. A non-zero
// status is returned for every create instead of a success.
func fakeCreateServer(t *testing.T, collections string, status int) (*[]createdLink, *atomic.Int64) {
	t.Helper()

	var (
		posts []createdLink
		count atomic.Int64
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/collections":
			_, _ = w.Write([]byte(`{"response":` + collections + `}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/search":
			_, _ = w.Write([]byte(`{"message":"Success","data":[]}`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/links":
			count.Add(1)
			body, _ := io.ReadAll(r.Body)
			var link createdLink
			require.NoError(t, json.Unmarshal(body, &link))
			posts = append(posts, link)
			if status != 0 {
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"response":"nope"}`))
				return
			}
			_, _ = w.Write([]byte(`{"response":{"id":42,"url":"https://go.dev"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	utils.SetVars("test-token", server.URL)

	return &posts, &count
}

const twoCollections = `[{"id":1,"name":"Unorganized"},{"id":7,"name":"Go"}]`

func TestCreateSavesSingleLink(t *testing.T) {
	assert := require.New(t)

	posts, _ := fakeCreateServer(t, twoCollections, 0)

	cmd := Create()
	cmd.SetArgs([]string{"https://go.dev", "--name", "Go", "--description", "The Go site",
		"--collection-id", "7", "--tag", "go", "--tag", "lang"})
	assert.NoError(cmd.Execute())

	assert.Len(*posts, 1)
	got := (*posts)[0]
	assert.Equal("https://go.dev", got.URL)
	assert.Equal("Go", got.Name)
	assert.Equal("The Go site", got.Description)
	assert.Equal(7, got.Collection.ID)
	assert.Len(got.Tags, 2)
	assert.Equal("go", got.Tags[0].Name)
	assert.Equal("lang", got.Tags[1].Name)
}

func TestCreateDefaultsToUnorganized(t *testing.T) {
	assert := require.New(t)

	posts, _ := fakeCreateServer(t, `[{"id":7,"name":"Go"},{"id":3,"name":"Unorganized"}]`, 0)

	cmd := Create()
	cmd.SetArgs([]string{"https://go.dev"})
	assert.NoError(cmd.Execute())

	assert.Len(*posts, 1)
	assert.Equal(3, (*posts)[0].Collection.ID)
}

func TestCreateUnknownCollectionErrors(t *testing.T) {
	assert := require.New(t)

	_, count := fakeCreateServer(t, twoCollections, 0)

	cmd := Create()
	cmd.SetArgs([]string{"https://go.dev", "--collection-id", "9999"})
	err := cmd.Execute()
	assert.ErrorIs(err, errCollectionNotFound)
	assert.Zero(count.Load())
}

func TestCreateConflictReportsAlreadyExists(t *testing.T) {
	assert := require.New(t)

	fakeCreateServer(t, twoCollections, http.StatusConflict)

	cmd := Create()
	cmd.SetArgs([]string{"https://go.dev"})
	err := cmd.Execute()
	assert.ErrorIs(err, errLinkExists)
	assert.ErrorContains(err, "https://go.dev")
}

func TestCreateRequiresURL(t *testing.T) {
	assert := require.New(t)

	for _, args := range [][]string{{}, {"https://a.example", "https://b.example"}} {
		cmd := Create()
		cmd.SetArgs(args)
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		assert.ErrorContains(cmd.Execute(), "requires the URL to save")
	}
}

func TestCreateRejectsFileFlag(t *testing.T) {
	assert := require.New(t)

	cmd := Create()
	cmd.SetArgs([]string{"--file", "links.json"})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	assert.ErrorContains(cmd.Execute(), "unknown flag: --file")
}
