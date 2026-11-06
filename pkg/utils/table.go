package utils

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog/log"
)

const (
	// ellipsis marks a cell value that was shortened to fit the terminal.
	ellipsis = "…"

	// defaultFlexMinWidth is the narrowest a flexible column is shrunk to
	// (unless its header is wider, or its content is narrower).
	defaultFlexMinWidth = 12
)

// Table collects rows for one of the CLI's result tables.
type Table struct {
	data    tui.Table
	columns []Column
	out     io.Writer
}

// Column describes one table column; see Table.Columns.
type Column struct {
	// Header is the column name.
	Header string
	// Flex marks a column whose values may be shortened with an ellipsis so
	// the printed table fits the terminal. Other columns keep their width.
	Flex bool
	// MinWidth is the narrowest a Flex column is shrunk to. Zero means
	// defaultFlexMinWidth. The header width is always respected.
	MinWidth int
	// HideWhenEmpty leaves the column out when it's empty in every row.
	HideWhenEmpty bool
}

// NewTable starts a table titled title, written to stdout.
func NewTable(title string) *Table {
	return &Table{data: tui.Table{Title: title}, out: os.Stdout}
}

// Header sets the column names. Every column keeps its full width; use
// Columns for columns that may be shortened or hidden.
func (t *Table) Header(cols ...string) {
	columns := make([]Column, len(cols))
	for i, c := range cols {
		columns[i] = Column{Header: c}
	}

	t.Columns(columns...)
}

// Columns sets the columns and how each one is fitted to the terminal.
func (t *Table) Columns(cols ...Column) {
	t.columns = cols
	t.data.Headers = make([]string, len(cols))
	for i, c := range cols {
		t.data.Headers[i] = c.Header
	}
}

// Row appends a row; each cell is formatted with fmt.Sprint. Line breaks and
// tabs become spaces so a value can't break the row across lines.
func (t *Table) Row(cells ...any) {
	row := make([]string, len(cells))
	for i, c := range cells {
		row[i] = singleLine(fmt.Sprint(c))
	}
	t.data.Rows = append(t.data.Rows, row)
}

// Tone colors data cells: fn gets each cell's 0-based row, column and text.
// Colors only show in static output when ColorEnabled.
func (t *Table) Tone(fn func(row, col int, value string) Tone) {
	t.data.Tone = func(row, col int) Tone {
		if row < 0 || row >= len(t.data.Rows) || col < 0 || col >= len(t.data.Rows[row]) {
			return ToneNone
		}
		return fn(row, col, t.data.Rows[row][col])
	}
}

// Tone is a semantic cell color; see Table.Tone.
type Tone = tui.Tone

const (
	ToneNone    = tui.ToneNone
	ToneSuccess = tui.ToneSuccess
	ToneWarning = tui.ToneWarning
	ToneDanger  = tui.ToneDanger
)

// Len is the number of data rows.
func (t *Table) Len() int { return len(t.data.Rows) }

// Render shows the table. When stdin and stdout are both terminals and the
// table is taller than the screen, it opens an interactive table (in the style
// of Bubble Tea's table example) that stays on screen after you quit.
// Otherwise (it fits, output is piped, or --non-interactive is set) it prints
// the table in the same style, in color only when ColorEnabled.
//
// When the printed table goes to a terminal, Flex columns are shortened with
// "…" so no line is wider than the terminal. Piped or redirected output is
// never shortened, so nothing is lost for grep or scripts. Either way,
// HideWhenEmpty columns that are empty in every row are left out.
func (t *Table) Render() {
	if width, height, ok := t.browseSize(); ok {
		err := tui.BrowseTable(t.visible(), width, height, tui.Options{})
		if err == nil {
			return
		}
		log.Debug().Err(err).Msg("table viewer failed, printing instead")
	}

	t.renderStatic(t.terminalWidth(), ColorEnabled())
}

// renderStatic prints the table, fitted to maxWidth columns; maxWidth <= 0
// means no limit.
func (t *Table) renderStatic(maxWidth int, color bool) {
	_ = tui.WriteStatic(t.out, fitTable(t.visible(), t.visibleColumns(), maxWidth), color)
}

// terminalWidth is the width to fit printed output to: the terminal's width
// when the table goes to a terminal, otherwise 0 (no limit).
func (t *Table) terminalWidth() int {
	if t.out != os.Stdout || !isatty.IsTerminal(os.Stdout.Fd()) {
		return 0
	}

	width, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 0
	}

	return width
}

// keptColumns is the index of every column that isn't hidden by
// HideWhenEmpty.
func (t *Table) keptColumns() []int {
	keep := make([]int, 0, len(t.columns))
	for i, c := range t.columns {
		if c.HideWhenEmpty && t.columnEmpty(i) {
			continue
		}

		keep = append(keep, i)
	}

	return keep
}

func (t *Table) columnEmpty(col int) bool {
	for _, row := range t.data.Rows {
		if col < len(row) && strings.TrimSpace(row[col]) != "" {
			return false
		}
	}

	return true
}

// visibleColumns is the Column of every column keptColumns keeps.
func (t *Table) visibleColumns() []Column {
	keep := t.keptColumns()
	cols := make([]Column, len(keep))
	for j, i := range keep {
		cols[j] = t.columns[i]
	}

	return cols
}

// visible is the table data without the columns keptColumns drops. Every
// row gets one cell per column, and Tone still sees the original columns.
func (t *Table) visible() tui.Table {
	keep := t.keptColumns()

	data := tui.Table{Title: t.data.Title, Headers: make([]string, len(keep))}
	for j, i := range keep {
		data.Headers[j] = t.data.Headers[i]
	}

	data.Rows = make([][]string, len(t.data.Rows))
	for r, row := range t.data.Rows {
		data.Rows[r] = make([]string, len(keep))
		for j, i := range keep {
			if i < len(row) {
				data.Rows[r][j] = row[i]
			}
		}
	}

	if tone := t.data.Tone; tone != nil {
		data.Tone = func(row, col int) Tone {
			if col < 0 || col >= len(keep) {
				return ToneNone
			}

			return tone(row, keep[col])
		}
	}

	return data
}

// fitTable shortens the Flex columns of data (described by cols) so the
// printed table is at most maxWidth wide. maxWidth <= 0 leaves it as is.
func fitTable(data tui.Table, cols []Column, maxWidth int) tui.Table {
	if maxWidth <= 0 {
		return data
	}

	natural := make([]int, len(cols))
	minimum := make([]int, len(cols))
	flex := make([]bool, len(cols))

	for i, col := range cols {
		headerWidth := ansi.StringWidth(col.Header)
		natural[i] = headerWidth
		for _, row := range data.Rows {
			natural[i] = max(natural[i], ansi.StringWidth(row[i]))
		}

		minWidth := col.MinWidth
		if minWidth <= 0 {
			minWidth = defaultFlexMinWidth
		}

		minimum[i] = min(natural[i], max(headerWidth, minWidth))
		flex[i] = col.Flex
	}

	widths := fitWidths(natural, minimum, flex, maxWidth-tui.StaticOverhead(len(cols)))

	rows := make([][]string, len(data.Rows))
	for r, row := range data.Rows {
		rows[r] = make([]string, len(row))
		for i, cell := range row {
			rows[r][i] = Truncate(cell, widths[i])
		}
	}

	data.Rows = rows

	return data
}

// fitWidths returns a width for each column so that their sum is at most
// budget. Non-flexible columns keep their natural width. Flexible columns
// share what is left: narrow ones keep their natural width and the rest split
// the remainder evenly, but none goes below its minimum. If even the minimums
// don't fit, the table ends up wider than budget; that's the best that can be
// done.
func fitWidths(natural, minimum []int, flex []bool, budget int) []int {
	widths := slices.Clone(natural)

	total := 0
	for _, w := range natural {
		total += w
	}

	if total <= budget {
		return widths
	}

	flexBudget := budget
	flexCols := make([]int, 0, len(natural))
	for i, w := range natural {
		if flex[i] {
			flexCols = append(flexCols, i)
		} else {
			flexBudget -= w
		}
	}

	// Water-fill: hand out the budget narrowest-first, so a column that
	// needs less than an even share gives the rest to the wider ones.
	slices.SortStableFunc(flexCols, func(a, b int) int { return cmp.Compare(natural[a], natural[b]) })

	remaining := flexBudget
	for k, i := range flexCols {
		share := remaining / (len(flexCols) - k)
		widths[i] = max(minimum[i], min(natural[i], share))
		remaining -= widths[i]
	}

	return widths
}

// Truncate shortens s to at most width terminal cells, ending it with "…"
// when anything was cut. Width is measured the way Lip Gloss lays tables out,
// so wide characters (emoji, CJK) count as two cells and are never split, and
// accented letters count as one. width <= 0 means no limit.
func Truncate(s string, width int) string {
	if width <= 0 || ansi.StringWidth(s) <= width {
		return s
	}

	cut := ansi.Truncate(s, width-ansi.StringWidth(ellipsis), "")

	return strings.TrimRight(cut, " ") + ellipsis
}

// singleLine replaces line breaks and tabs with spaces.
func singleLine(s string) string {
	return strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ", "\t", " ").Replace(s)
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
