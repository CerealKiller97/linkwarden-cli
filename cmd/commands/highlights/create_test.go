package highlights

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestCreate_UsesLinkIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	link := e2etest.Link(t, client, e2etest.Collection(t, client).ID)

	cmd := Create()
	cmd.SetArgs([]string{
		"--link-id", strconv.Itoa(link.ID),
		"--text", "quoted",
		"--color", "yellow",
		"--start-offset", "0",
		"--end-offset", "6",
	})
	assert.NoError(cmd.Execute())

	highlights, err := client.Links.Highlights(t.Context(), link.ID)
	assert.NoError(err)
	assert.Len(highlights, 1)
	assert.Equal("quoted", highlights[0].Text)
	assert.Equal("yellow", highlights[0].Color)
	assert.Equal(6, highlights[0].EndOffset)
}

func TestCreate_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := Create()
	cmd.SetArgs([]string{
		"--link-id", "9",
		"--text", "quoted",
		"--color", "yellow",
		"--start-offset", "0",
		"--end-offset", "6",
	})
	assert.Error(cmd.Execute())
}
