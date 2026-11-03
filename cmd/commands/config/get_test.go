package config

import (
	"os"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGet_RendersConfig(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Get()
	assert.NoError(cmd.Execute())
}

func TestGet_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	// The config endpoint is public, so point at a path the instance doesn't serve instead.
	e2etest.Setup(t)
	utils.SetVars(os.Getenv("LINKWARDEN_TOKEN"), os.Getenv("LINKWARDEN_URL")+"/no-such-prefix")

	cmd := Get()
	assert.Error(cmd.Execute())
}
