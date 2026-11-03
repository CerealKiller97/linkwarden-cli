package favicon

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestGet_SavesToOutputFile(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)
	outputPath := filepath.Join(t.TempDir(), "favicon")

	cmd := Get()
	cmd.SetArgs([]string{"--url", "https://go.dev", "--output", outputPath})
	assert.NoError(cmd.Execute())

	data, err := os.ReadFile(outputPath)
	assert.NoError(err)
	assert.NotEmpty(data)
}

func TestGet_WithoutOutputPrintsSummary(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	cmd := Get()
	cmd.SetArgs([]string{"--url", "https://go.dev"})
	assert.NoError(cmd.Execute())
}

func TestGet_PropagatesAPIError(t *testing.T) {
	assert := require.New(t)

	// The favicon endpoint is public, so point at a path the instance doesn't serve instead.
	e2etest.Setup(t)
	utils.SetVars(os.Getenv("LINKWARDEN_TOKEN"), os.Getenv("LINKWARDEN_URL")+"/no-such-prefix")

	cmd := Get()
	cmd.SetArgs([]string{"--url", "https://go.dev"})
	assert.Error(cmd.Execute())
}
