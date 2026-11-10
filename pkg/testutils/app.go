// Package testutils holds helpers for tests. NewApp builds the *app.App a
// test runs a command with, so each test has its own credentials, output
// mode, output buffer and prompter instead of sharing package variables, and
// can run in parallel.
package testutils

import (
	"bytes"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
)

// Token is the API token every test App carries.
const Token = "test-token"

// WithPrompter answers prompts with p (such as a prompttest.Script), and
// treats stdin as an interactive terminal (--non-interactive still wins).
func WithPrompter(p app.Prompter) app.Option {
	return func(app *app.App) {
		app.Prompter = p
		app.Scripted = true
	}
}

// NewApp returns an App that talks to instanceURL with Token and prints
// results to the returned buffer. opts (app.WithMode, app.WithNonInteractive,
// WithPrompter, ...) apply on top. Pass the App to the command under test.
func NewApp(instanceURL string, opts ...app.Option) (*app.App, *bytes.Buffer) {
	out := &bytes.Buffer{}
	base := []app.Option{app.WithCredentials(Token, instanceURL), app.WithOutput(out)}

	return app.New(append(base, opts...)...), out
}
