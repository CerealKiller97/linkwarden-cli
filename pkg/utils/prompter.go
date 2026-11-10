package utils

import "github.com/CerealKiller97/linkwarden-cli/pkg/app"

// Prompter answers the CLI's single-question prompts; see app.Prompter.
type Prompter = app.Prompter

// The legacy prompter, used only when a command runs without an *app.App in
// its context. Removed once every command and test runs with one.
var (
	prompter Prompter = app.TerminalPrompter{}
	// scriptedPrompter is set while a test's Prompter is in place: there's
	// someone to answer, so prompts run even though stdin isn't a terminal.
	scriptedPrompter bool
)

// SetPrompter makes p answer prompts until restore is called, and treats stdin
// as an interactive terminal meanwhile (--non-interactive still wins). It only
// affects commands run without an *app.App in their context.
//
// Legacy, removed once every command reads the App: run the command with testutils.WithPrompter instead.
func SetPrompter(p Prompter) (restore func()) {
	prev, prevScripted := prompter, scriptedPrompter
	prompter, scriptedPrompter = p, true

	return func() {
		prompter, scriptedPrompter = prev, prevScripted
	}
}
