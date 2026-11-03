package archives

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

// linkWithPNGArchive creates a link with png as its uploaded PNG archive.
func linkWithPNGArchive(t *testing.T, c *linkwarden.Client, png []byte) *linkwarden.Link {
	t.Helper()

	link := e2etest.Link(t, c, e2etest.Collection(t, c).ID)
	_, err := c.Archives.Upload(t.Context(), link.ID, linkwarden.ArchiveFormatPNG, false, "screenshot.png", bytes.NewReader(png))
	require.NoError(t, err)

	return link
}

func TestGet_SavesToOutputFile(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	png := e2etest.PNG(t)
	link := linkWithPNGArchive(t, client, png)
	outputPath := filepath.Join(t.TempDir(), "archive.png")

	cmd := Get()
	cmd.SetArgs([]string{"--link-id", strconv.Itoa(link.ID), "--format", "0", "--output", outputPath})
	assert.NoError(cmd.Execute())

	data, err := os.ReadFile(outputPath)
	assert.NoError(err)
	assert.Equal(png, data)
}

func TestGet_WithoutOutputPrintsSummary(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	link := linkWithPNGArchive(t, client, e2etest.PNG(t))

	cmd := Get()
	cmd.SetArgs([]string{"--link-id", strconv.Itoa(link.ID)})
	assert.NoError(cmd.Execute())
}

func TestGet_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Get()
	cmd.SetArgs([]string{"--link-id", "999999999"})
	assert.Error(cmd.Execute())
}
