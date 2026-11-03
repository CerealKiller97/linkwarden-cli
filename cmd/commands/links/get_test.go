package links

import (
	"strconv"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGet_UsesIDFlag(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	link := e2etest.Link(t, client, e2etest.Collection(t, client).ID)

	cmd := Get()
	cmd.SetArgs([]string{"--id", strconv.Itoa(link.ID)})
	assert.NoError(cmd.Execute())
}

func TestGet_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Get()
	cmd.SetArgs([]string{"--id", "999999999"})
	assert.Error(cmd.Execute())
}

func TestGet_InvalidIDErrors(t *testing.T) {
	assert := require.New(t)

	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Get()
	cmd.SetArgs([]string{"--id", "not-a-number"})
	assert.Error(cmd.Execute())
}
