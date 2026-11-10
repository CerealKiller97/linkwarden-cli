package testutils_test

import (
	"bytes"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewAppAppliesOptions(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	const instanceURL = "http://example.invalid"

	// Act
	ctx, a := testutils.NewApp(t, instanceURL, app.WithMode(app.OutputJSON), app.WithNonInteractive(true))

	// Assert
	assert.Same(a, app.From(ctx))
	assert.Equal(testutils.Token, a.Token)
	assert.Equal(instanceURL, a.InstanceURL)
	assert.Equal(app.OutputJSON, a.Mode)
	assert.True(a.NonInteractive)
	assert.IsType(&bytes.Buffer{}, a.Out)
}

func TestNewAppContextReturnsOutputBuffer(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	ctx, out := testutils.NewAppContext(t, "http://example.invalid")

	// Act
	_, err := app.From(ctx).Stdout().Write([]byte("hello"))

	// Assert
	assert.NoError(err)
	assert.Equal("hello", out.String())
}

func TestWithPrompterMarksPromptsScripted(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	p := app.TerminalPrompter{}

	// Act
	_, a := testutils.NewApp(t, "http://example.invalid", testutils.WithPrompter(p))

	// Assert
	assert.Equal(p, a.Prompts())
	assert.True(a.Scripted)
}
