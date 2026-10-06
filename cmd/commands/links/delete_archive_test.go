package links

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestDeleteArchive_ParsesIDs(t *testing.T) {
	assert := require.New(t)

	var gotBody struct {
		LinkIDs []int `json:"linkIds"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/links/archive", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Archived files deleted."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := DeleteArchive()
	cmd.SetArgs([]string{"--ids", "1, 2,3"})
	assert.NoError(cmd.Execute())
	assert.Equal([]int{1, 2, 3}, gotBody.LinkIDs)
}

func TestDeleteArchive_InvalidIDsErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := DeleteArchive()
	cmd.SetArgs([]string{"--ids", "not-a-number"})
	assert.Error(cmd.Execute())
}
