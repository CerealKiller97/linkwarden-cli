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
