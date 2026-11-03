package migrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestExport_WritesOutputFile(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	collection := e2etest.Collection(t, client)
	link := e2etest.Link(t, client, collection.ID)
	outputPath := filepath.Join(t.TempDir(), "export.json")

	cmd := Export()
	cmd.SetArgs([]string{"--output", outputPath})
	assert.NoError(cmd.Execute())

	data, err := os.ReadFile(outputPath)
	assert.NoError(err)
	assert.Contains(string(data), collection.Name)
	assert.Contains(string(data), link.URL)
}

func TestExport_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := Export()
	cmd.SetArgs([]string{"--output", filepath.Join(t.TempDir(), "export.json")})
	assert.Error(cmd.Execute())
}
