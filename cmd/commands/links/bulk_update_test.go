package links

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/linkwarden"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestBulkUpdate_BuildsRequestFromFlags(t *testing.T) {
	assert := require.New(t)

	var gotBody linkwarden.BulkUpdateLinksRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/links", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Links updated."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := BulkUpdate()
	cmd.SetArgs([]string{
		"--ids", "1,2",
		"--collection-id", "7",
		"--tag", "a",
		"--tag", "b",
		"--remove-previous-tags",
	})
	assert.NoError(cmd.Execute())

	assert.Len(gotBody.Links, 2)
	assert.Equal(7, gotBody.NewData.CollectionID)
	assert.Equal([]linkwarden.Tag{{Name: "a"}, {Name: "b"}}, gotBody.NewData.Tags)
	assert.True(gotBody.RemovePreviousTags)
}

func TestBulkUpdate_InvalidIDsErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := BulkUpdate()
	cmd.SetArgs([]string{"--ids", "not-a-number"})
	assert.Error(cmd.Execute())
}
