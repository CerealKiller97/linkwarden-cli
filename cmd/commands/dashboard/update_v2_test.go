package dashboard

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestUpdateV2_AppliesLayoutFromFile(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)
	filePath := e2etest.WriteFile(t, "layout.json", []byte(
		`[{"type":"STATS","enabled":true,"order":0},{"type":"RECENT_LINKS","enabled":true,"order":1}]`))

	cmd := UpdateV2()
	cmd.SetArgs([]string{"--file", filePath})
	assert.NoError(cmd.Execute())
}

func TestUpdateV2_MissingFileErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := UpdateV2()
	cmd.SetArgs([]string{"--file", "/nonexistent/layout.json"})
	assert.Error(cmd.Execute())
}
