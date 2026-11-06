package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/mattn/go-isatty"
)

// PickerItem is a single choice in an interactive ID picker.
type PickerItem struct {
	Label string
	Value int
}

// SelectID shows an interactive promptui select over items and returns the chosen
// item's Value.
func SelectID(label string, items []PickerItem) (int, error) {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Label
	}

	prompt := promptui.Select{
		Label: label,
		Items: names,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	return items[idx].Value, nil
}

// ResolveID returns flagValue parsed as an int when non-empty. Otherwise, if stdin
// is an interactive terminal, it calls fetch to list the available items and shows
// a picker for label. If stdin isn't a terminal (a script, CI, piped input, ...),
// it fails fast with an actionable error instead of hanging on a prompt nobody can
// answer. flagHint (e.g. "--id") names the flag to use for non-interactive callers,
// and is only used in that error message.
func ResolveID(flagValue, flagHint, label string, fetch func() ([]PickerItem, error)) (int, error) {
	if flagValue != "" {
		return strconv.Atoi(flagValue)
	}

	if !isInteractive() {
		return 0, fmt.Errorf("no %s given and no interactive terminal available; pass %s explicitly for non-interactive use", flagHint, flagHint)
	}

	items, err := fetch()
	if err != nil {
		return 0, err
	}

	return SelectID(label, items)
}

func isInteractive() bool {
	return isInteractiveFor(GetNonInteractive(), isatty.IsTerminal(os.Stdin.Fd()))
}

// isInteractiveFor is the pure decision behind isInteractive: the
// --non-interactive override always wins; otherwise it's whatever stdin
// actually is.
func isInteractiveFor(forcedNonInteractive, isTerminal bool) bool {
	if forcedNonInteractive {
		return false
	}

	return isTerminal
}

// Select shows a promptui select over options and returns the chosen index and
// option. Without an interactive terminal it fails fast instead of hanging.
func Select(label string, options []string) (int, string, error) {
	if !isInteractive() {
		return 0, "", errNoTerminal(label)
	}

	prompt := promptui.Select{Label: label, Items: options}

	return prompt.Run()
}

// Input asks for a line of text. def, when non-empty, is pre-filled and used if
// the user just presses Enter. validate (if non-nil) must pass before the prompt
// accepts the value. Without an interactive terminal it fails fast.
func Input(label, def string, validate func(string) error) (string, error) {
	if !isInteractive() {
		return "", errNoTerminal(label)
	}

	prompt := promptui.Prompt{Label: label, Default: def, AllowEdit: def != "", Validate: validate}

	return prompt.Run()
}

// Password is Input with the typed characters masked and no default.
func Password(label string, validate func(string) error) (string, error) {
	if !isInteractive() {
		return "", errNoTerminal(label)
	}

	prompt := promptui.Prompt{Label: label, Mask: '*', Validate: validate}

	return prompt.Run()
}

func errNoTerminal(label string) error {
	return fmt.Errorf("%q needs an interactive terminal, but none is available", label)
}
