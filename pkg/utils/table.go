package utils

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

const (
	// ellipsis marks a cell value that was shortened to fit the terminal.
	ellipsis = "…"

	// defaultFlexMinWidth is the narrowest a flexible column is shrunk to
	// (unless its header is wider, or its content is narrower).
	defaultFlexMinWidth = 12
)

// NewTable returns a table.Writer preconfigured with the CLI's shared look
// (rounded borders, bold cyan header) writing to stdout. Color is skipped
// automatically when stdout isn't a terminal or NO_COLOR is set, so piped or
// redirected output stays free of raw ANSI escape codes.
//
// NewTable doesn't fit its content to the terminal; use NewListTable for
// tables whose cells can be long (names, URLs, descriptions).
func NewTable(title string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(title)
	t.SetStyle(tableStyleFor(ColorEnabled()))

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

// Column describes one column of a ListTable.
type Column struct {
	// Header is the column title (rendered upper-case).
	Header string
	// Flex marks a column whose values may be shortened with an ellipsis so
	// the table fits the terminal. Columns that aren't flexible keep their
	// full width.
	Flex bool
	// MinWidth is the narrowest a Flex column is shrunk to. Zero means
	// defaultFlexMinWidth. The header width is always respected.
	MinWidth int
	// HideWhenEmpty drops the column when it's empty in every row.
	HideWhenEmpty bool
}

// ListTable is a table of rows that is sized to fit the terminal when it's
// rendered: on a terminal, flexible columns are shortened with an ellipsis
// so no line is wider than the terminal. When stdout isn't a terminal (piped
// or redirected), values are never shortened, so nothing is lost for tools
// like grep.
type ListTable struct {
	title   string
	columns []Column
	rows    [][]string
}

// NewListTable returns an empty ListTable with the given title and columns.
func NewListTable(title string, columns ...Column) *ListTable {
	return &ListTable{title: title, columns: columns}
}

// AppendRow adds a row. Cells are formatted with fmt.Sprint; line breaks and
// tabs are turned into spaces so every row stays on one line. Missing cells
// are left empty and extra cells are ignored.
func (lt *ListTable) AppendRow(cells ...any) {
	row := make([]string, len(lt.columns))
	for i := range row {
		if i < len(cells) {
			row[i] = singleLine(fmt.Sprint(cells[i]))
		}
	}

	lt.rows = append(lt.rows, row)
}

// Render writes the table to stdout, fitted to the terminal width when
// stdout is a terminal.
func (lt *ListTable) Render() {
	lt.RenderTo(os.Stdout, TerminalWidth(), ColorEnabled())
}

// RenderTo writes the table to w. maxWidth is the widest a line may be;
// zero or less means no limit (values are never shortened).
func (lt *ListTable) RenderTo(w io.Writer, maxWidth int, colorEnabled bool) {
	columns, rows := lt.visibleColumns()

	natural := make([]int, len(columns))
	minimum := make([]int, len(columns))
	flex := make([]bool, len(columns))

	for i, col := range columns {
		headerWidth := text.StringWidthWithoutEscSequences(strings.ToUpper(col.Header))
		natural[i] = headerWidth

		for _, row := range rows {
			natural[i] = max(natural[i], text.StringWidthWithoutEscSequences(row[i]))
		}

		flex[i] = col.Flex
		minWidth := col.MinWidth
		if minWidth <= 0 {
			minWidth = defaultFlexMinWidth
		}

		minimum[i] = min(natural[i], max(headerWidth, minWidth))
	}

	widths := natural
	if maxWidth > 0 {
		// The title row spans the cells and the borders between them.
		titleCells := text.StringWidthWithoutEscSequences(lt.title) - 3*(len(columns)-1)
		widths = fitWidths(natural, minimum, flex, maxWidth-tableOverhead(len(columns)), titleCells)
	}

	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.SetTitle(lt.title)
	t.SetStyle(tableStyleFor(colorEnabled))

	header := make(table.Row, len(columns))
	for i, col := range columns {
		header[i] = col.Header
	}

	t.AppendHeader(header)

	for _, row := range rows {
		out := make(table.Row, len(row))
		for i, cell := range row {
			out[i] = Truncate(cell, widths[i])
		}

		t.AppendRow(out)
	}

	t.Render()
}

// visibleColumns returns the columns (and the matching cells of each row)
// left after dropping HideWhenEmpty columns that are empty in every row.
func (lt *ListTable) visibleColumns() ([]Column, [][]string) {
	keep := make([]int, 0, len(lt.columns))

	for i, col := range lt.columns {
		if col.HideWhenEmpty && lt.columnEmpty(i) {
			continue
		}

		keep = append(keep, i)
	}

	columns := make([]Column, len(keep))
	for j, i := range keep {
		columns[j] = lt.columns[i]
	}

	rows := make([][]string, len(lt.rows))
	for r, row := range lt.rows {
		rows[r] = make([]string, len(keep))
		for j, i := range keep {
			rows[r][j] = row[i]
		}
	}

	return columns, rows
}

func (lt *ListTable) columnEmpty(i int) bool {
	for _, row := range lt.rows {
		if strings.TrimSpace(row[i]) != "" {
			return false
		}
	}

	return true
}

// tableOverhead is the number of columns the rounded style uses for borders
// and padding around n cells: "│ " before each cell, " " after it, and one
// closing "│".
func tableOverhead(n int) int {
	return 3*n + 1
}

// fitWidths returns a width for each column so that their sum is at most
// budget (and at least minTotal, so a title spanning the table still fits).
// Non-flexible columns keep their natural width. Flexible columns share what
// is left: narrow ones keep their natural width and the rest split the
// remainder evenly, but none goes below its minimum. If even the minimums
// don't fit, the minimums are returned (the table will be wider than
// budget, which is the best that can be done).
func fitWidths(natural, minimum []int, flex []bool, budget, minTotal int) []int {
	widths := append([]int(nil), natural...)

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

	if len(flexCols) == 0 {
		return widths
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

	// Widen the last (widest) flexible column if the title wouldn't fit.
	sum := 0
	for _, w := range widths {
		sum += w
	}

	if sum < minTotal {
		last := flexCols[len(flexCols)-1]
		widths[last] += minTotal - sum
	}

	return widths
}

// Truncate shortens s to at most width terminal columns, ending it with "…"
// when anything was cut. Width is measured in terminal cells, so wide
// characters (emoji, CJK) count as two and combining marks as zero; a wide
// character that would straddle the limit is dropped rather than split.
// width <= 0 means no limit.
func Truncate(s string, width int) string {
	if width <= 0 || text.StringWidthWithoutEscSequences(s) <= width {
		return s
	}

	limit := width - text.StringWidthWithoutEscSequences(ellipsis)

	var b strings.Builder

	used := 0
	for _, r := range s {
		w := text.RuneWidth(r)
		if used+w > limit {
			break
		}

		b.WriteRune(r)
		used += w
	}

	return strings.TrimRight(b.String(), " ") + ellipsis
}

// singleLine replaces line breaks and tabs with spaces so a value can't
// break a table row across lines.
func singleLine(s string) string {
	return strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ", "\t", " ").Replace(s)
}

// TerminalWidth returns the width of the terminal stdout is attached to, or
// 0 when stdout isn't a terminal (piped or redirected) or its size can't be
// read. 0 means "no limit" to ListTable.
func TerminalWidth() int {
	fd := os.Stdout.Fd()
	if !isatty.IsTerminal(fd) {
		return 0
	}

	width, _, err := term.GetSize(int(fd)) //nolint:gosec // a file descriptor always fits in an int
	if err != nil {
		return 0
	}

	return width
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
