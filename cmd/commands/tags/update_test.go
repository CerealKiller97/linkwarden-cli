package tags

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestUpdateUsesIDAndNameFlags(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v1/tags/3", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":3,"name":"Renamed"}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update(nil)
	cmd.SetArgs([]string{"--id", "3", "--name", "Renamed"})
	assert.NoError(cmd.Execute())
}

func TestUpdatePropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found"}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Update(nil)
	cmd.SetArgs([]string{"--id", "3", "--name", "Renamed"})
	assert.Error(cmd.Execute())
}
