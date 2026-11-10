package testutils_test

import (
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
	wantMode := app.OutputJSON
	opts := []app.Option{app.WithMode(wantMode), app.WithNonInteractive(true)}

	// Act
	app, out := testutils.NewApp(instanceURL, opts...)

	// Assert
	assert.Equal(testutils.Token, app.Token)
	assert.Equal(instanceURL, app.InstanceURL)
	assert.Equal(wantMode, app.Mode)
	assert.True(app.NonInteractive)
	assert.Same(out, app.Stdout())
}

func TestNewAppPrintsToReturnedBuffer(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	app, out := testutils.NewApp("http://example.invalid")

	// Act
	_, err := app.Stdout().Write([]byte("hello"))

	// Assert
	assert.NoError(err)
	assert.Equal("hello", out.String())
}

func TestWithPrompterMarksPromptsScripted(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	prompter := app.TerminalPrompter{}

	// Act
	app, _ := testutils.NewApp("http://example.invalid", testutils.WithPrompter(prompter))

	// Assert
	assert.Equal(prompter, app.Prompts())
	assert.True(app.Scripted)
}
