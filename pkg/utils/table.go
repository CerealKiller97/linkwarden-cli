package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-isatty"
)

// NewTable returns a table.Writer preconfigured with the CLI's shared look
// (rounded borders, bold cyan header) writing to stdout. Color is skipped
// automatically when stdout isn't a terminal or NO_COLOR is set, so piped or
// redirected output stays free of raw ANSI escape codes.
func NewTable(title string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(title)
	t.SetStyle(tableStyleFor(colorEnabled()))

	return t
}

// tableStyleFor returns the shared rounded-border style, adding a bold cyan
// header (foreground only, no forced background) when colorEnabled is true.
func tableStyleFor(colorEnabled bool) table.Style {
	style := table.StyleRounded
	if !colorEnabled {
		return style
	}

	style.Color.Header = text.Colors{text.FgHiCyan, text.Bold}

	return style
}

func colorEnabled() bool {
	return colorEnabledFor(os.Getenv("NO_COLOR"), isatty.IsTerminal(os.Stdout.Fd()))
}

// colorEnabledFor is the pure decision behind colorEnabled: color is on only
// when NO_COLOR isn't set and output is going to an actual terminal.
func colorEnabledFor(noColorEnv string, isTerminal bool) bool {
	if noColorEnv != "" {
		return false
	}

	return isTerminal
}
