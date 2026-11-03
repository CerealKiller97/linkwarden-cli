package dashboard

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestGet_RendersLinks(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	e2etest.Link(t, client, e2etest.Collection(t, client).ID)

	cmd := Get()
	assert.NoError(cmd.Execute())
}

func TestGet_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := Get()
	assert.Error(cmd.Execute())
}
