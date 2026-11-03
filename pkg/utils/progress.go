package utils

import (
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/mattn/go-isatty"
)

// Progress is a progress bar for a known amount of work.
type Progress interface {
	Increment()
	Stop()
}

// StartProgress draws a Bubble Tea progress bar on stderr when stderr is a
// terminal; otherwise it returns a no-op, so logs and CI output stay clean.
func StartProgress(label string, total int) Progress {
	if !isatty.IsTerminal(os.Stderr.Fd()) {
		return noProgress{}
	}

	return tui.StartProgress(label, total, os.Stderr)
}

type noProgress struct{}

func (noProgress) Increment() {}
func (noProgress) Stop()      {}
