package dashboard

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestGetV2_RendersStatsAndLinks(t *testing.T) {
	assert := require.New(t)

	client := e2etest.Setup(t)
	e2etest.Link(t, client, e2etest.Collection(t, client).ID)

	cmd := GetV2()
	assert.NoError(cmd.Execute())
}

func TestGetV2_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := GetV2()
	assert.Error(cmd.Execute())
}
