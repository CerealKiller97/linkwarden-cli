// Package testutils holds helpers for tests. NewAppContext and NewApp build the *app.App
// a test runs a command with, so each test has its own credentials, output
// mode, output buffer and prompter instead of sharing package variables, and
// can run in parallel.
package testutils

import (
	"bytes"
	"context"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
)

// Token is the API token every test App carries.
const Token = "test-token"

// WithPrompter answers prompts with p (such as a prompttest.Script), and
// treats stdin as an interactive terminal (--non-interactive still wins).
func WithPrompter(p app.Prompter) app.Option {
	return func(a *app.App) {
		a.Prompter = p
		a.Scripted = true
	}
}

// NewAppContext returns a context carrying an App that talks to instanceURL with Token
// and prints results to the returned buffer. opts (app.WithMode,
// app.WithNonInteractive, WithPrompter, ...) apply on top. Run the command with
// cmd.ExecuteContext(ctx).
func NewAppContext(t testing.TB, instanceURL string, opts ...app.Option) (context.Context, *bytes.Buffer) {
	t.Helper()

	ctx, a := NewApp(t, instanceURL, opts...)

	return ctx, a.Out.(*bytes.Buffer)
}

// NewApp is NewAppContext for tests that need the App itself, for example to read its
// InstanceURL.
func NewApp(t testing.TB, instanceURL string, opts ...app.Option) (context.Context, *app.App) {
	t.Helper()

	base := []app.Option{app.WithCredentials(Token, instanceURL), app.WithOutput(&bytes.Buffer{})}
	a := app.New(append(base, opts...)...)

	return app.WithApp(t.Context(), a), a
}
