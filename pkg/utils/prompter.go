package utils

import "github.com/CerealKiller97/linkwarden-cli/pkg/tui"

// Prompter answers the CLI's single-question prompts: Select, Input,
// InputDefault and Password, and so Confirm too. The default one draws them in
// the terminal with pkg/tui; tests swap in scripted answers with SetPrompter
// (see pkg/utils/prompttest). The paged ID pickers behind ResolveID and
// SelectID aren't routed through it; tests pass the ID flag instead.
type Prompter interface {
	Select(label string, options []string) (int, error)
	Input(label string, validate func(string) error) (string, error)
	InputDefault(label, def string, validate func(string) error) (string, error)
	Password(label string, validate func(string) error) (string, error)
}

// terminalPrompter is the default Prompter: Bubble Tea prompts in the terminal.
type terminalPrompter struct{}

func (terminalPrompter) Select(label string, options []string) (int, error) {
	return tui.Select(label, options)
}

func (terminalPrompter) Input(label string, validate func(string) error) (string, error) {
	return tui.Input(label, validate)
}

func (terminalPrompter) InputDefault(label, def string, validate func(string) error) (string, error) {
	return tui.InputDefault(label, def, validate)
}

func (terminalPrompter) Password(label string, validate func(string) error) (string, error) {
	return tui.Password(label, validate)
}

var (
	prompter Prompter = terminalPrompter{}
	// scriptedPrompter is set while a test's Prompter is in place: there's
	// someone to answer, so prompts run even though stdin isn't a terminal.
	scriptedPrompter bool
)

// SetPrompter makes p answer prompts until restore is called, and treats stdin
// as an interactive terminal meanwhile (--non-interactive still wins). It is
// meant for tests, which can't type into a real terminal.
func SetPrompter(p Prompter) (restore func()) {
	prev, prevScripted := prompter, scriptedPrompter
	prompter, scriptedPrompter = p, true

	return func() {
		prompter, scriptedPrompter = prev, prevScripted
	}
}
