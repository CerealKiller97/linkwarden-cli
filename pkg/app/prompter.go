package app

import "github.com/CerealKiller97/linkwarden-cli/pkg/tui"

// Prompter answers the CLI's single-question prompts: Select, Input,
// InputDefault and Password, and so Confirm too. The default one draws them in
// the terminal with pkg/tui; tests set App.Prompter to scripted answers (see
// pkg/utils/prompttest). The paged ID pickers behind ResolveID and SelectID
// aren't routed through it; tests pass the ID flag instead.
type Prompter interface {
	Select(label string, options []string) (int, error)
	Input(label string, validate func(string) error) (string, error)
	InputDefault(label, def string, validate func(string) error) (string, error)
	Password(label string, validate func(string) error) (string, error)
}

// TerminalPrompter is the default Prompter: Bubble Tea prompts in the terminal.
type TerminalPrompter struct{}

func (TerminalPrompter) Select(label string, options []string) (int, error) {
	return tui.Select(label, options)
}

func (TerminalPrompter) Input(label string, validate func(string) error) (string, error) {
	return tui.Input(label, validate)
}

func (TerminalPrompter) InputDefault(label, def string, validate func(string) error) (string, error) {
	return tui.InputDefault(label, def, validate)
}

func (TerminalPrompter) Password(label string, validate func(string) error) (string, error) {
	return tui.Password(label, validate)
}
