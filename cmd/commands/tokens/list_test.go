package tokens

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestList_RendersTokens(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := List()
	assert.NoError(cmd.Execute())
}

func TestList_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := List()
	assert.Error(cmd.Execute())
}
