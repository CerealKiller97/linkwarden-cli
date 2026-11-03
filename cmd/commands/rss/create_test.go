package rss

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestCreate_SendsRequestedFields(t *testing.T) {
	assert := require.New(t)

	var gotBody linkwarden.CreateRSSRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPost, r.Method)
		assert.Equal("/api/v1/rss", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":1,"name":"News","url":"https://news.example.com/rss"}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Create()
	cmd.SetArgs([]string{"--name", "News", "--url", "https://news.example.com/rss", "--collection-id", "7"})
	assert.NoError(cmd.Execute())
	assert.Equal("News", gotBody.Name)
	assert.Equal(7, gotBody.CollectionID)
}

func TestCreate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"invalid feed"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Create()
	cmd.SetArgs([]string{"--name", "News", "--url", "https://news.example.com/rss"})
	assert.Error(cmd.Execute())
}
