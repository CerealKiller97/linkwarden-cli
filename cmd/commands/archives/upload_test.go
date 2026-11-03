package archives

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/stretchr/testify/require"
)

func TestUpload_SendsFile(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	link := e2etest.Link(t, client, e2etest.Collection(t, client).ID)
	png := e2etest.PNG(t)
	filePath := e2etest.WriteFile(t, "screenshot.png", png)

	cmd := Upload()
	cmd.SetArgs([]string{"--link-id", strconv.Itoa(link.ID), "--file", filePath, "--format", "0"})
	assert.NoError(cmd.Execute())

	file, err := client.Archives.Get(t.Context(), link.ID, linkwarden.ArchiveFormatPNG, false)
	assert.NoError(err)
	assert.Equal(png, file.Data)
}

func TestUpload_MissingFileErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Upload()
	cmd.SetArgs([]string{"--link-id", "2", "--file", "/nonexistent/file.png", "--format", "0"})
	assert.Error(cmd.Execute())
}
