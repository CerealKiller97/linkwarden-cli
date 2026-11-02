package linkwarden

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPublicService_CollectionLinks(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/public/collections/links", r.URL.Path)
		q := r.URL.Query()
		assert.Equal("4", q.Get("collectionId"))
		assert.Equal("go", q.Get("searchQueryString"))
		assert.Equal("true", q.Get("searchByUrl"))
		assert.False(q.Has("searchByName"))
		assert.False(q.Has("cursor"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":1,"name":"Go","url":"https://go.dev"}]}`))
	})

	links, err := client.Public.CollectionLinks(t.Context(), 4, ListPublicLinksOptions{SearchQueryString: "go", SearchByURL: true})
	assert.NoError(err)
	assert.Len(links, 1)
}

func TestPublicService_CollectionTags(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal("/api/v1/public/collections/tags", r.URL.Path)
		assert.Equal("4", r.URL.Query().Get("collectionId"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":[{"id":2,"name":"Technology"}]}`))
	})

	tags, err := client.Public.CollectionTags(t.Context(), 4)
	assert.NoError(err)
	assert.Equal("Technology", tags[0].Name)
}

func TestPublicService_GetByID(t *testing.T) {
	assert := require.New(t)

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		switch r.URL.Path {
		case "/api/v1/public/collections/3":
			_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Reading"}}`))
		case "/api/v1/public/links/5":
			_, _ = w.Write([]byte(`{"response":{"id":5,"name":"Go","url":"https://go.dev"}}`))
		case "/api/v1/public/users/1":
			_, _ = w.Write([]byte(`{"response":{"id":1,"username":"dev","archiveAsPDF":true}}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
	})

	collection, err := client.Public.Collection(t.Context(), 3)
	assert.NoError(err)
	assert.Equal("Reading", collection.Name)

	link, err := client.Public.Link(t.Context(), 5)
	assert.NoError(err)
	assert.Equal("Go", link.Name)

	user, err := client.Public.User(t.Context(), 1)
	assert.NoError(err)
	assert.Equal("dev", user.Username)
	assert.True(user.ArchiveAsPDF)
}
