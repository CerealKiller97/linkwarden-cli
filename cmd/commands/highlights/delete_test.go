package highlights

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestDelete_UsesIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	link := e2etest.Link(t, client, e2etest.Collection(t, client).ID)
	highlight := e2etest.Highlight(t, client, link.ID)

	cmd := Delete()
	cmd.SetArgs([]string{"--id", strconv.Itoa(highlight.ID)})
	assert.NoError(cmd.Execute())

	highlights, err := client.Links.Highlights(t.Context(), link.ID)
	assert.NoError(err)
	assert.Empty(highlights)
}

func TestDelete_InvalidIDErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Delete()
	cmd.SetArgs([]string{"--id", "not-a-number"})
	assert.Error(cmd.Execute())
}
