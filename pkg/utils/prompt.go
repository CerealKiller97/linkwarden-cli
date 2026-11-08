package utils

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/mattn/go-isatty"
)

// PickerItem is a single choice in an interactive ID picker. Detail, if set, is
// shown dimmed after the label (and can be filtered on), e.g. a link's URL.
type PickerItem struct {
	Label  string
	Detail string
	Value  int
}

// PageFetch fetches one page of picker items starting at cursor (0 for the
// first page), and the cursor of the page after it, or 0 when there is none.
type PageFetch func(cursor int) (items []PickerItem, next int, err error)

// SelectID shows an interactive picker over items and returns the chosen item's
// Value.
func SelectID(label string, items []PickerItem) (int, error) {
	idx, err := tui.SelectPaged(label, selectOptions(items), nil, tui.Options{})
	if err != nil {
		return 0, err
	}

	return items[idx].Value, nil
}

func selectOptions(items []PickerItem) []tui.SelectOption {
	opts := make([]tui.SelectOption, len(items))
	for i, item := range items {
		opts[i] = tui.SelectOption{Label: item.Label, Detail: item.Detail}
	}

	return opts
}

// ResolveIDPaged is ResolveID for long lists the API returns a page at a
// time: the picker shows the first page and fetches the next ones as the user
// pages or filters past what's loaded.
func ResolveIDPaged(flagValue, flagHint, label string, fetch PageFetch) (int, error) {
	if flagValue != "" {
		return strconv.Atoi(flagValue)
	}

	if !isInteractive() {
		return 0, fmt.Errorf("no %s given and no interactive terminal available; pass %s explicitly for non-interactive use", flagHint, flagHint)
	}

	spinner := StartSpinner("Loading choices…")
	items, next, err := fetch(0)
	spinner.Stop()
	if err != nil {
		return 0, err
	}

	loaded := &pickerPages{items: items, next: next, fetch: fetch}

	idx, err := tui.SelectPaged(label, selectOptions(items), loaded.loader(), tui.Options{})
	if err != nil {
		return 0, err
	}

	return loaded.value(idx), nil
}

// pickerPages holds every page of picker items loaded so far, in the order
// the picker shows them, so its indexes line up. Pages load on a Bubble Tea
// goroutine, possibly still running when the user picks, hence the lock.
type pickerPages struct {
	mu    sync.Mutex
	items []PickerItem
	next  int
	fetch PageFetch
}

// loader returns the tui.LoadMore that fetches the next page, or nil when the
// first page was the only one.
func (p *pickerPages) loader() tui.LoadMore {
	if p.next == 0 {
		return nil
	}

	return func() ([]tui.SelectOption, bool, error) {
		p.mu.Lock()
		cursor := p.next
		p.mu.Unlock()

		page, next, err := p.fetch(cursor)
		if err != nil {
			return nil, false, err
		}

		p.mu.Lock()
		p.items = append(p.items, page...)
		p.next = next
		p.mu.Unlock()

		return selectOptions(page), next != 0, nil
	}
}

// value is the Value of the idx-th item loaded.
func (p *pickerPages) value(idx int) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.items[idx].Value
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
	return ResolveIDPaged(flagValue, flagHint, label, func(int) ([]PickerItem, int, error) {
		items, err := fetch()
		return items, 0, err
	})
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
