package links

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestGetHighlights_UsesIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	link := e2etest.Link(t, client, e2etest.Collection(t, client).ID)
	e2etest.Highlight(t, client, link.ID)

	cmd := GetHighlights()
	cmd.SetArgs([]string{"--id", strconv.Itoa(link.ID)})
	assert.NoError(cmd.Execute())
}

func TestGetHighlights_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := GetHighlights()
	cmd.SetArgs([]string{"--id", "6"})
	assert.Error(cmd.Execute())
}
