package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func writeLayoutFile(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "layout.json")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	return path
}

func TestUpdateV2SendsParsedLayout(t *testing.T) {
	assert := require.New(t)

	var gotBody []linkwarden.DashboardLayoutItem
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(http.MethodPut, r.Method)
		assert.Equal("/api/v2/dashboard", r.URL.Path)
		assert.NoError(json.NewDecoder(r.Body).Decode(&gotBody))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true,"message":"Success","data":{"links":[],"numberOfPinnedLinks":0,"numberOfTags":0}}`))
	}))
	defer server.Close()
	utils.SetVars("test-token", server.URL)

	filePath := writeLayoutFile(t, `[{"type":"STATS","enabled":true,"order":0}]`)

	cmd := UpdateV2()
	cmd.SetArgs([]string{"--file", filePath})
	assert.NoError(cmd.Execute())
	assert.Len(gotBody, 1)
	assert.Equal("STATS", gotBody[0].Type)
}

func TestUpdateV2MissingFileErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := UpdateV2()
	cmd.SetArgs([]string{"--file", "/nonexistent/layout.json"})
	assert.Error(cmd.Execute())
}
