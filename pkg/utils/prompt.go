package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/mattn/go-isatty"
)

// PickerItem is a single choice in an interactive ID picker.
type PickerItem struct {
	Label string
	Value int
}

// SelectID shows an interactive picker over items and returns the chosen item's
// Value.
func SelectID(label string, items []PickerItem) (int, error) {
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Label
	}

	idx, err := tui.Select(label, names)
	if err != nil {
		return 0, err
	}

	return items[idx].Value, nil
}

// Select shows an interactive picker over options and returns the chosen index
// and option. Unlike ResolveID there's no flag to fall back on, so without an
// interactive terminal it fails fast instead of waiting on a prompt.
func Select(label string, options []string) (int, string, error) {
	if !isInteractive() {
		return 0, "", errNoTerminal(label)
	}

	idx, err := tui.Select(label, options)
	if err != nil {
		return 0, "", err
	}

	return idx, options[idx], nil
}

// Input asks for a line of text, re-prompting until validate (if non-nil)
// passes. Without an interactive terminal it fails fast.
func Input(label string, validate func(string) error) (string, error) {
	if !isInteractive() {
		return "", errNoTerminal(label)
	}

	return tui.Input(label, validate)
}

func errNoTerminal(label string) error {
	return fmt.Errorf("%q needs an interactive terminal, but none is available", label)
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

// InputDefault is Input with def pre-filled as a placeholder: pressing Enter on
// an empty field accepts def. Without an interactive terminal it fails fast.
func InputDefault(label, def string, validate func(string) error) (string, error) {
	if !isInteractive() {
		return "", errNoTerminal(label)
	}

	return tui.InputDefault(label, def, validate)
}

// Password is Input with the typed characters masked. Without an interactive
// terminal it fails fast.
func Password(label string, validate func(string) error) (string, error) {
	if !isInteractive() {
		return "", errNoTerminal(label)
	}

	return tui.Password(label, validate)
}
