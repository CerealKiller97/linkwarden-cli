package tags

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
	tag := e2etest.Tag(t, client)

	cmd := Get()
	cmd.SetArgs([]string{"--id", strconv.Itoa(tag.ID)})
	assert.NoError(cmd.Execute())
}

func TestGet_NotFoundErrors(t *testing.T) {
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
