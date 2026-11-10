// Package apptest builds the *app.App a test runs a command with, so each test
// has its own credentials, output mode, output buffer and prompter instead of
// sharing package variables, and can run in parallel.
package apptest

import (
	"bytes"
	"context"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
)

// Token is the API token every test App carries.
const Token = "test-token"

// Option changes the App New builds.
type Option func(*app.App)

// WithMode prints results in mode instead of tables.
func WithMode(mode app.OutputMode) Option {
	return func(a *app.App) { a.Mode = mode }
}

// NonInteractive behaves as if --non-interactive was passed.
func NonInteractive() Option {
	return func(a *app.App) { a.NonInteractive = true }
}

// WithPrompter answers prompts with p (such as a prompttest.Script), and
// treats stdin as an interactive terminal (--non-interactive still wins).
func WithPrompter(p app.Prompter) Option {
	return func(a *app.App) {
		a.Prompter = p
		a.Scripted = true
	}
}

// New returns a context carrying an App that talks to instanceURL with Token
// and prints results to the returned buffer. Run the command with
// cmd.ExecuteContext(ctx).
func New(t testing.TB, instanceURL string, opts ...Option) (context.Context, *bytes.Buffer) {
	t.Helper()

	ctx, a := NewApp(t, instanceURL, opts...)

	return ctx, a.Out.(*bytes.Buffer)
}

// NewApp is New for tests that need the App itself, for example to read its
// InstanceURL.
func NewApp(t testing.TB, instanceURL string, opts ...Option) (context.Context, *app.App) {
	t.Helper()

	a := &app.App{Token: Token, InstanceURL: instanceURL, Out: &bytes.Buffer{}}
	for _, opt := range opts {
		opt(a)
	}

	return app.WithApp(t.Context(), a), a
}
