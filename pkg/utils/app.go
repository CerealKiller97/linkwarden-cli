package utils

import (
	"context"

	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
)

// current returns the App the helpers in this package work with: the one in
// ctx, or, for commands and tests that don't carry one yet, an App built from
// the legacy package variables (SetVars, SetOutputMode, SetNonInteractive,
// SetPrompter). An App in ctx without a Prompter of its own also falls back to
// the legacy prompter, so prompttest.Use keeps working under the root command.
func current(ctx context.Context) *app.App {
	if a, ok := app.FromContext(ctx); ok {
		if a.Prompter != nil {
			return a
		}

		withPrompter := *a
		withPrompter.Prompter, withPrompter.Scripted = prompter, scriptedPrompter

		return &withPrompter
	}

	return &app.App{
		Token:          LinkwardenToken,
		InstanceURL:    InstanceURL,
		Mode:           outputMode,
		NonInteractive: NonInteractive,
		Prompter:       prompter,
		Scripted:       scriptedPrompter,
	}
}
