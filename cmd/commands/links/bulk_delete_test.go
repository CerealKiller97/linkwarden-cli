package links

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
		LinkIDs []int `json:"linkIds"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodDelete, r.Method)
		assert.Equal("/api/v1/links", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":"Links deleted."}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", "4,5"})
	assert.NoError(cmd.Execute())
	assert.Equal([]int{4, 5}, gotBody.LinkIDs)
}

func TestBulkDeletePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"boom"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := BulkDelete()
	cmd.SetArgs([]string{"--ids", "4"})
	assert.Error(cmd.Execute())
}
