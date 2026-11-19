package utils

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/mattn/go-isatty"
)

// Spinner shows that something is in progress, for work of unknown length
// such as an API request.
type Spinner interface {
	Stop()
}

// StartSpinner draws a Bubble Tea spinner labelled label on stderr when stderr
// is a terminal; otherwise it returns a no-op, so logs and CI output stay clean.
// Call Stop as soon as the work is done, before printing results or opening a
// prompt: it erases the spinner so the next output starts on a clean line.
func StartSpinner(label string) Spinner {
	if !isatty.IsTerminal(os.Stderr.Fd()) {
		return noSpinner{}
	}

	return tui.StartSpinner(label, os.Stderr)
}

type noSpinner struct{}

func (noSpinner) Stop() {}
