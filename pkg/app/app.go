// Package app holds the state one run of the CLI works with: the instance and
// token to talk to, how to print results, and whether prompts may be shown.
// The root command builds an App from its flags, the environment and the
// config file, and stores it in the command's context; commands and the
// helpers in pkg/utils read it back with From.
//
// Keeping this out of package variables lets tests build their own App per
// test (see pkg/testutils) and run in parallel.
package app

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/CerealKiller97/linkwarden-go-sdk"
	"github.com/rs/zerolog/log"
)

// OutputMode is how commands print their results on stdout. It is set from the
// global --json and --quiet flags.
type OutputMode int

const (
	// OutputTable prints the CLI's usual tables (the default).
	OutputTable OutputMode = iota
	// OutputJSON prints results as indented JSON, using the API's field names.
	OutputJSON
	// OutputIDs prints only the ID of each result, one per line, for piping
	// into another command. Tables without an ID column print nothing.
	OutputIDs
)

// App is the state of one run of the CLI. Build one with New. The zero value
// is usable too: table output on stdout, terminal prompts, no credentials.
type App struct {
	// Out is where results are printed. Nil means os.Stdout.
	Out io.Writer
	// Prompter answers prompts. Nil means TerminalPrompter.
	Prompter Prompter

	// Token and InstanceURL are the credentials API requests use.
	Token       string
	InstanceURL string

	// Mode is how results are printed; see OutputMode.
	Mode OutputMode

	// NonInteractive forces non-interactive mode (--non-interactive), so
	// commands fail fast instead of prompting.
	NonInteractive bool
	// Scripted reports that Prompter answers from a script (tests), so prompts
	// run even though stdin isn't a terminal. --non-interactive still wins.
	Scripted bool
}

// Option configures the App New builds.
type Option func(*App)

// WithCredentials sets the token and instance URL API requests use.
func WithCredentials(token, instanceURL string) Option {
	return func(a *App) {
		a.Token = token
		a.InstanceURL = instanceURL
	}
}

// WithMode sets how results are printed; see OutputMode.
func WithMode(mode OutputMode) Option {
	return func(a *App) { a.Mode = mode }
}

// WithNonInteractive forces non-interactive mode when on (--non-interactive).
func WithNonInteractive(on bool) Option {
	return func(a *App) { a.NonInteractive = on }
}

// WithOutput prints results to w instead of os.Stdout.
func WithOutput(w io.Writer) Option {
	return func(a *App) { a.Out = w }
}

// WithPrompter answers prompts with p instead of the terminal.
func WithPrompter(p Prompter) Option {
	return func(a *App) { a.Prompter = p }
}

// New returns an App that prints tables to os.Stdout and prompts in the
// terminal, with no credentials, changed by opts.
func New(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Client returns a Linkwarden API client for the App's instance and token.
func (a *App) Client() *linkwarden.Client {
	return linkwarden.New(a.InstanceURL, a.Token, linkwarden.WithLogger(slogBridge()))
}

// Stdout is where results are printed: Out, or os.Stdout when it's nil.
func (a *App) Stdout() io.Writer {
	if a.Out == nil {
		return os.Stdout
	}

	return a.Out
}

// Prompts returns the Prompter that answers prompts: Prompter, or the terminal
// when it's nil.
func (a *App) Prompts() Prompter {
	if a.Prompter == nil {
		return TerminalPrompter{}
	}

	return a.Prompter
}

type ctxKey struct{}

// WithApp returns a copy of ctx that carries a.
func WithApp(ctx context.Context, a *App) context.Context {
	return context.WithValue(ctx, ctxKey{}, a)
}

// FromContext returns the App stored in ctx with WithApp, if any.
func FromContext(ctx context.Context) (*App, bool) {
	if ctx == nil {
		return nil, false
	}

	a, ok := ctx.Value(ctxKey{}).(*App)

	return a, ok && a != nil
}

// From returns the App stored in ctx with WithApp. Every command runs under
// the root command, which stores one, and tests store their own; a missing App
// is a bug, so From panics rather than guess.
func From(ctx context.Context) *App {
	a, ok := FromContext(ctx)
	if !ok {
		panic("app: no *app.App in the command's context; run it under the root command or use testutils.NewApp")
	}

	return a
}

// slogBridge adapts the global zerolog logger to the *slog.Logger the SDK
// expects. zerolog.Logger is an io.Writer, so slog's text output is forwarded
// through it; slog's own timestamp is dropped because zerolog adds one.
func slogBridge() *slog.Logger {
	return slog.New(slog.NewTextHandler(log.Logger, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
}
