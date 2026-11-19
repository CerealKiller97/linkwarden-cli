package utils

import (
	"github.com/CerealKiller97/linkwarden-cli/pkg/app"
)

// orLegacy returns the App the helpers in this package work with: app, or,
// when it's nil (commands and tests that don't pass one yet), an App built
// from the legacy package variables (SetVars, SetOutputMode,
// SetNonInteractive, SetPrompter). An App without a Prompter of its own also
// falls back to the legacy prompter, so prompttest.Use keeps working under
// the root command.
func orLegacy(app *app.App) *app.App {
	if app == nil {
		return legacyApp()
	}

	if app.Prompter != nil {
		return app
	}

	withPrompter := *app
	withPrompter.Prompter, withPrompter.Scripted = prompter, scriptedPrompter

	return &withPrompter
}

// legacyApp builds an App from the legacy package variables.
func legacyApp() *app.App {
	return &app.App{
		Token:          LinkwardenToken,
		InstanceURL:    InstanceURL,
		Mode:           outputMode,
		NonInteractive: NonInteractive,
		Prompter:       prompter,
		Scripted:       scriptedPrompter,
	}
}
