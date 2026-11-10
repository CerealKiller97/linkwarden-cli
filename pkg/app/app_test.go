package app_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/stretchr/testify/require"
)

func TestFromReturnsStoredApp(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	a := app.New(app.WithCredentials("t", "https://links.example.com"))
	ctx := app.WithApp(t.Context(), a)

	// Act
	got := app.From(ctx)

	// Assert
	assert.Same(a, got)
}

func TestFromPanicsWithoutApp(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	ctx := context.Background()

	// Act & Assert
	assert.Panics(func() { app.From(ctx) })
}

func TestNewAppliesOptions(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	var out bytes.Buffer
	p := app.TerminalPrompter{}

	// Act
	a := app.New(
		app.WithCredentials("tok", "https://links.example.com"),
		app.WithMode(app.OutputIDs),
		app.WithNonInteractive(true),
		app.WithOutput(&out),
		app.WithPrompter(p),
	)

	// Assert
	assert.Equal("tok", a.Token)
	assert.Equal("https://links.example.com", a.InstanceURL)
	assert.Equal(app.OutputIDs, a.Mode)
	assert.True(a.NonInteractive)
	assert.Same(&out, a.Stdout())
	assert.Equal(p, a.Prompts())
	assert.False(a.Scripted)
}

func TestNewDefaults(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Act
	a := app.New()

	// Assert
	assert.Equal(os.Stdout, a.Stdout())
	assert.IsType(app.TerminalPrompter{}, a.Prompts())
	assert.Equal(app.OutputTable, a.Mode)
	assert.False(a.NonInteractive)
	assert.Empty(a.Token)
}

func TestZeroAppDefaults(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	// Arrange
	var a app.App

	// Act
	out, prompter := a.Stdout(), a.Prompts()

	// Assert
	assert.Equal(os.Stdout, out)
	assert.IsType(app.TerminalPrompter{}, prompter)
	assert.Equal(app.OutputTable, a.Mode)
}
