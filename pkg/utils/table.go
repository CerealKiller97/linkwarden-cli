package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/charmbracelet/x/term"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog/log"
)

// Table collects rows for one of the CLI's result tables.
type Table struct {
	data tui.Table
	out  io.Writer
}

// NewTable starts a table titled title, written to stdout.
func NewTable(title string) *Table {
	return &Table{data: tui.Table{Title: title}, out: os.Stdout}
}

// Header sets the column names.
func (t *Table) Header(cols ...string) { t.data.Headers = cols }

// Row appends a row; each cell is formatted with fmt.Sprint.
func (t *Table) Row(cells ...any) {
	row := make([]string, len(cells))
	for i, c := range cells {
		row[i] = fmt.Sprint(c)
	}
	t.data.Rows = append(t.data.Rows, row)
}

// Len is the number of data rows.
func (t *Table) Len() int { return len(t.data.Rows) }

// Render shows the table. When stdin and stdout are both terminals and the
// table is taller than the screen, it opens an interactive table (in the style
// of Bubble Tea's table example) that stays on screen after you quit.
// Otherwise (it fits, output is piped, or --non-interactive is set) it prints
// the table in the same style, in color only when ColorEnabled.
func (t *Table) Render() {
	if width, height, ok := t.browseSize(); ok {
		err := tui.BrowseTable(t.data, width, height, tui.Options{})
		if err == nil {
			return
		}
		log.Debug().Err(err).Msg("table viewer failed, printing instead")
	}

	_ = tui.WriteStatic(t.out, t.data, ColorEnabled())
}

// browseSize returns the terminal size when the viewer should be used.
func (t *Table) browseSize() (width, height int, ok bool) {
	if t.out != os.Stdout || !isInteractive() || !isatty.IsTerminal(os.Stdout.Fd()) {
		return 0, 0, false
	}

	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 0, 0, false
	}

	return width, height, shouldBrowse(len(t.data.Rows), height)
}

// shouldBrowse reports whether rows data rows (plus the table's borders,
// header and title, and a line for the next prompt) overflow a terminal
// height rows tall.
func shouldBrowse(rows, height int) bool {
	return height > 0 && rows+tui.TableChromeHeight+1 > height
}

// ColorEnabled reports whether CLI output should include ANSI color: only
// when NO_COLOR isn't set and stdout is an actual terminal (not piped or
// redirected).
func ColorEnabled() bool {
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
