package tags

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestList_RendersTags(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	tag := e2etest.Tag(t, client)

	cmd := List()
	cmd.SetArgs([]string{"--search", tag.Name})
	assert.NoError(cmd.Execute())
}

func TestList_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := List()
	assert.Error(cmd.Execute())
}
