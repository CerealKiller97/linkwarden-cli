package tags

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestBulkDeleteParsesIDs(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		TagIDs []int `json:"tagIds"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/tags", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Tags deleted."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", "1,2"})
	assert.NoError(cmd.Execute())
	assert.Equal([]int{1, 2}, gotBody.TagIDs)
}

func TestBulkDeleteInvalidIDsErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", "not-a-number"})
	assert.Error(cmd.Execute())
}
