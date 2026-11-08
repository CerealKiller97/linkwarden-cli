package links

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGetUsesIDFlag(t *testing.T) {
	assert := require.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodGet, r.Method)
		assert.Equal("/api/v1/links/5", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"response":{"id":5,"name":"Example","url":"https://example.com"}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	cmd := Get()
	cmd.SetArgs([]string{"--id", "5"})
	assert.NoError(cmd.Execute())
}

func TestGetInvalidIDErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Get()
	cmd.SetArgs([]string{"--id", "not-a-number"})
	assert.Error(cmd.Execute())
}
