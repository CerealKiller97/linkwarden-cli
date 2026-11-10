package app_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
	"github.com/CerealKiller97/linkwarden-cli/pkg/app/apptest"
	"github.com/stretchr/testify/require"
)

func TestFromReturnsStoredApp(t *testing.T) {
	t.Parallel()

	a := app.New(app.WithCredentials("t", "https://links.example.com"))
	require.Same(t, a, app.From(app.WithApp(t.Context(), a)))
}

func TestFromPanicsWithoutApp(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() { app.From(context.Background()) })
}

func TestNewAppliesOptions(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	p := app.TerminalPrompter{}
	a := app.New(
		app.WithCredentials("tok", "https://links.example.com"),
		app.WithMode(app.OutputIDs),
		app.WithNonInteractive(true),
		app.WithOutput(&out),
		app.WithPrompter(p),
	)
	require.Equal(t, "tok", a.Token)
	require.Equal(t, "https://links.example.com", a.InstanceURL)
	require.Equal(t, app.OutputIDs, a.Mode)
	require.True(t, a.NonInteractive)
	require.Same(t, &out, a.Stdout())
	require.Equal(t, p, a.Prompts())
	require.False(t, a.Scripted)
}

func TestNewDefaults(t *testing.T) {
	t.Parallel()

	a := app.New()
	require.Equal(t, os.Stdout, a.Stdout())
	require.IsType(t, app.TerminalPrompter{}, a.Prompts())
	require.Equal(t, app.OutputTable, a.Mode)
	require.False(t, a.NonInteractive)
	require.Empty(t, a.Token)
}

func TestZeroAppDefaults(t *testing.T) {
	t.Parallel()

	var a app.App
	require.Equal(t, os.Stdout, a.Stdout())
	require.IsType(t, app.TerminalPrompter{}, a.Prompts())
	require.Equal(t, app.OutputTable, a.Mode)
}

func TestApptestNewAppliesOptions(t *testing.T) {
	t.Parallel()

	ctx, a := apptest.NewApp(t, "http://example.invalid", app.WithMode(app.OutputJSON), app.WithNonInteractive(true))
	require.Same(t, a, app.From(ctx))
	require.Equal(t, apptest.Token, a.Token)
	require.Equal(t, "http://example.invalid", a.InstanceURL)
	require.Equal(t, app.OutputJSON, a.Mode)
	require.True(t, a.NonInteractive)
	require.IsType(t, &bytes.Buffer{}, a.Out)
}
