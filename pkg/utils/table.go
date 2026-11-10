package utils

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/charmbracelet/x/term"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog/log"
)

// Table collects rows for one of the CLI's result tables.
type Table struct {
	out    io.Writer
	source any
	mode   *OutputMode // nil: the global mode; tests set it
	data   tui.Table
}

// JSON sets the value --json prints instead of the table: the API's own
// records (e.g. []linkwarden.Link, or one *linkwarden.Link for a get), so
// scripts see every field, not just the columns shown. Without it, --json
// prints the rows as objects keyed by header.
func (t *Table) JSON(v any) *Table {
	t.source = v
	return t
}

// Every table in the CLI is built with NewTable (a list of records, one per
// row, with the row count in the title) or NewSummary (a fixed set of rows,
// such as settings or totals, without a count). Both share one look: Title
// Case headers, ID first when there is one, and cells formatted by FormatCell.

// NewTable starts a list table titled title (plural, e.g. "Links"), written to
// stdout. Its title shows the row count: "Links (12)".
func NewTable(title string) *Table {
	return &Table{data: tui.Table{Title: title, ShowCount: true}, out: os.Stdout}
}

// NewSummary starts a table whose rows aren't a list of records, like a
// key/value listing or a row of totals, so its title has no row count.
func NewSummary(title string) *Table {
	return &Table{data: tui.Table{Title: title}, out: os.Stdout}
}

// Caption sets a dim line shown between the title and the table, for totals
// or other context that isn't a row.
func (t *Table) Caption(text string) { t.data.Caption = text }

// Header sets the column names.
func (t *Table) Header(cols ...string) { t.data.Headers = cols }

// Row appends a row; each cell is formatted with FormatCell.
func (t *Table) Row(cells ...any) {
	row := make([]string, len(cells))
	for i, c := range cells {
		row[i] = FormatCell(c)
	}
	t.data.Rows = append(t.data.Rows, row)
}

// EmptyCell is what a table shows for a missing or blank value.
const EmptyCell = "-"

// cellTimeLayout is how table cells show a time, in the local time zone.
const cellTimeLayout = "2006-01-02 15:04"

// FormatCell turns a value into table text the same way everywhere: nil, nil
// pointers, blank strings and zero times become EmptyCell, booleans become
// Yes/No, times use cellTimeLayout, and pointers show what they point to.
// Anything else uses fmt.Sprint.
func FormatCell(v any) string {
	switch v := v.(type) {
	case nil:
		return EmptyCell
	case string:
		if strings.TrimSpace(v) == "" {
			return EmptyCell
		}
		return v
	case bool:
		if v {
			return "Yes"
		}
		return "No"
	case time.Time:
		if v.IsZero() {
			return EmptyCell
		}
		return v.Local().Format(cellTimeLayout)
	}

	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return EmptyCell
		}
		return FormatCell(rv.Elem().Interface())
	}

	return fmt.Sprint(v)
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

// LinkColumn makes the http(s) URLs in column col (0-based) clickable in
// terminals that support hyperlinks. A long URL is still cut short with an
// ellipsis to fit, but the link opens the whole URL. Cells that aren't an
// http(s) URL, such as EmptyCell, stay plain text.
func (t *Table) LinkColumn(col int) {
	t.data.Link = func(row, c int) string {
		if c != col || row < 0 || row >= len(t.data.Rows) || col >= len(t.data.Rows[row]) {
			return ""
		}

		return linkTarget(t.data.Rows[row][col])
	}
}

// linkTarget returns s when it's an absolute http or https URL, else "".
func linkTarget(s string) string {
	u, err := url.Parse(strings.TrimSpace(s))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return ""
	}

	return u.String()
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
// With --json it prints the value set with JSON instead, and with --quiet only
// the ID column.
func (t *Table) Render() {
	mode := outputMode
	if t.mode != nil {
		mode = *t.mode
	}

	switch mode {
	case OutputJSON:
		if err := WriteJSON(t.out, t.jsonValue()); err != nil {
			log.Error().Err(err).Msg("write JSON")
		}
		return
	case OutputIDs:
		t.writeIDs()
		return
	case OutputTable:
	}

	if width, height, ok := t.browseSize(); ok {
		err := tui.BrowseTable(t.data, width, height, tui.Options{})
		if err == nil {
			return
		}
		log.Debug().Err(err).Msg("table viewer failed, printing instead")
	}

	data := t.data
	if t.out == os.Stdout && isatty.IsTerminal(os.Stdout.Fd()) {
		if width, _, err := term.GetSize(os.Stdout.Fd()); err == nil {
			data.Width = width
		}
	}

	_ = tui.WriteStatic(t.out, data, ColorEnabled())
}

// jsonValue is what --json prints: the source set with JSON, or else the rows
// as objects keyed by header.
func (t *Table) jsonValue() any {
	if t.source != nil {
		return t.source
	}

	rows := make([]map[string]string, len(t.data.Rows))
	for i, row := range t.data.Rows {
		obj := make(map[string]string, len(row))
		for j, cell := range row {
			if j < len(t.data.Headers) {
				obj[t.data.Headers[j]] = cell
			}
		}
		rows[i] = obj
	}

	return rows
}

// writeIDs prints the ID column, one value per line. A table without an ID
// column (a summary) prints nothing.
func (t *Table) writeIDs() {
	col := -1
	for i, h := range t.data.Headers {
		if h == "ID" {
			col = i
			break
		}
	}
	if col < 0 {
		return
	}

	for _, row := range t.data.Rows {
		if col < len(row) {
			_, _ = fmt.Fprintln(t.out, row[col])
		}
	}
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

	return width, height, shouldBrowse(len(t.data.Rows), tui.ChromeHeight(t.data), height)
}

// shouldBrowse reports whether rows data rows (plus chrome rows of borders,
// header, title and caption, and a line for the next prompt) overflow a
// terminal height rows tall.
func shouldBrowse(rows, chrome, height int) bool {
	return height > 0 && rows+chrome+1 > height
}

// ColorSwatch is a table cell for a #rrggbb color: a dot in that color when
// output is colored (see ColorEnabled), otherwise the hex code itself, so
// piped output keeps the value. Anything that isn't a hex color is shown as is.
func ColorSwatch(hex string) string {
	if ColorEnabled() {
		if dot, ok := tui.Swatch(hex); ok {
			return dot
		}
	}

	return FormatCell(hex)
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
