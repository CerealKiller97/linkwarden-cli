package app_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/stretchr/testify/require"
)

func TestNewAppliesOptions(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	var out bytes.Buffer
	prompter := app.TerminalPrompter{}
	wantMode := app.OutputIDs
	opts := []app.Option{
		app.WithCredentials("tok", "https://links.example.com"),
		app.WithMode(wantMode),
		app.WithNonInteractive(true),
		app.WithOutput(&out),
		app.WithPrompter(prompter),
	}

	// Act
	app := app.New(opts...)

	// Assert
	assert.Equal("tok", app.Token)
	assert.Equal("https://links.example.com", app.InstanceURL)
	assert.Equal(wantMode, app.Mode)
	assert.True(app.NonInteractive)
	assert.Same(&out, app.Stdout())
	assert.Equal(prompter, app.Prompts())
	assert.False(app.Scripted)
}

func TestNewDefaults(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	wantMode := app.OutputTable
	wantPrompter := app.TerminalPrompter{}

	// Act
	app := app.New()

	// Assert
	assert.Equal(os.Stdout, app.Stdout())
	assert.IsType(wantPrompter, app.Prompts())
	assert.Equal(wantMode, app.Mode)
	assert.False(app.NonInteractive)
	assert.Empty(app.Token)
}

func TestZeroAppDefaults(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	wantMode := app.OutputTable
	wantPrompter := app.TerminalPrompter{}
	var app app.App

	// Act
	out, prompter := app.Stdout(), app.Prompts()

	// Assert
	assert.Equal(os.Stdout, out)
	assert.IsType(wantPrompter, prompter)
	assert.Equal(wantMode, app.Mode)
}
