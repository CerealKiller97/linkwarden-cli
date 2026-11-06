package tui

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	lgtable "charm.land/lipgloss/v2/table"
	"github.com/charmbracelet/x/ansi"
)

// Table is plain tabular data: a title, column headers and string cells.
type Table struct {
	Title   string
	Headers []string
	Rows    [][]string

	// ShowCount appends the number of rows to the title, as in "Links (12)".
	// List tables set it; summaries and key/value tables don't.
	ShowCount bool

	// Caption, when set, is a dim line between the title and the table, for
	// totals or other context that isn't a row.
	Caption string

	// Width, when positive, is the terminal width static output must fit in;
	// long cells are cut short with an ellipsis. Zero means no limit.
	Width int

	// Tone, when set, picks a color for a data cell (row and col are 0-based
	// into Rows). Only static output uses it.
	Tone func(row, col int) Tone
}

// Tone is a semantic color for a table cell.
type Tone int

const (
	ToneNone Tone = iota
	ToneSuccess
	ToneWarning
	ToneDanger
)

// Tables use the Linkwarden blues from theme.go: a rounded Ocean border, Sky
// header, a Sky title badge like Linkwarden's buttons, and muted alternate rows.
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Ink).
			Background(Sky).
			Padding(0, 1)
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Ocean)
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
	headerCell  = cellStyle.Bold(true).Foreground(Sky)
	borderStyle = lipgloss.NewStyle().Foreground(Ocean)

	toneStyles = map[Tone]lipgloss.Style{
		ToneSuccess: lipgloss.NewStyle().Bold(true).Foreground(Sky),
		ToneWarning: lipgloss.NewStyle().Bold(true).Foreground(Warning),
		ToneDanger:  lipgloss.NewStyle().Bold(true).Foreground(Danger),
	}
)

// RenderStatic returns t as text in the same style as the viewer. With color
// false the output has no ANSI escapes at all, so piped output is clean.
func RenderStatic(t Table, color bool) string {
	t = fitCells(t)
	tbl := lgtable.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(borderStyle).
		BorderColumn(false).
		BorderHeader(true).
		Headers(t.Headers...).
		Rows(t.Rows...).
		StyleFunc(func(row, col int) lipgloss.Style { return staticCellStyle(t, row, col) })

	// Plain output gets a bare title: the badge's padding only makes sense
	// with its background color behind it.
	if !color {
		return ansi.Strip(header(t, false) + tbl.String() + "\n")
	}

	return header(t, true) + tbl.String() + "\n"
}

func staticCellStyle(t Table, row, col int) lipgloss.Style {
	if row == lgtable.HeaderRow {
		return headerCell
	}

	if t.Tone != nil {
		if style, ok := toneStyles[t.Tone(row, col)]; ok {
			return style.Padding(0, 1)
		}
	}

	if row%2 == 1 {
		return cellStyle.Foreground(Muted)
	}

	return cellStyle
}

// header is the title line (with the row count when t.ShowCount) and the
// caption line, if any, that go above the table.
func header(t Table, color bool) string {
	title := t.Title
	if t.ShowCount {
		title = fmt.Sprintf("%s (%d)", title, len(t.Rows))
	}
	if color {
		title = titleStyle.Render(title)
	}

	out := title + "\n"
	if t.Caption != "" {
		out += dimStyle.Render(t.Caption) + "\n"
	}

	return out
}

// TableChromeHeight is how many terminal rows RenderStatic adds around the data
// rows of a table without a caption: title, top border, header, header rule
// and bottom border.
const TableChromeHeight = 5

// ChromeHeight is how many terminal rows RenderStatic adds around t's data
// rows: TableChromeHeight plus the caption line, if t has one.
func ChromeHeight(t Table) int {
	if t.Caption != "" {
		return TableChromeHeight + 1
	}

	return TableChromeHeight
}

// BrowseTable shows t as an interactive bubbles/table, inline in the terminal,
// fitted to width×height (the terminal's size). It returns when the user quits
// with q, Esc or Ctrl+C, leaving the table on screen.
func BrowseTable(t Table, width, height int, o Options) error {
	_, err := run(newBrowseModel(t, width, height), o)
	if errors.Is(err, ErrAborted) {
		return nil
	}

	return err
}

type browseModel struct {
	header string
	// extra is how many rows the header takes beyond the title line.
	extra int
	table table.Model
	done  bool
}

// browseChrome is the rows the viewer needs besides data rows: title, top
// border, header, header rule, bottom border, help line, and one spare line
// so the shell prompt doesn't scroll the table.
const browseChrome = 7

func newBrowseModel(t Table, width, height int) browseModel {
	// The outer border takes two columns.
	cols := fitColumns(t, width-2)
	extra := ChromeHeight(t) - TableChromeHeight
	rows := make([]table.Row, len(t.Rows))
	for i, r := range t.Rows {
		row := make(table.Row, len(r))
		for j, c := range r {
			row[j] = cleanCell(c)
		}
		rows[i] = row
	}

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Ocean).
		BorderBottom(true).
		Bold(true).
		Foreground(Sky)
	// Like the selected row in Linkwarden's sidebar.
	styles.Selected = styles.Selected.
		Foreground(Sky).
		Background(Slate).
		Bold(true)

	tbl := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		// Styles before height: WithHeight measures the (styled) header.
		table.WithStyles(styles),
		table.WithHeight(tableHeight(len(rows), height-extra)),
		// Without a width the viewport is zero columns wide and shows no rows.
		table.WithWidth(tableWidth(cols)),
	)

	return browseModel{header: header(t, true), extra: extra, table: tbl}
}

// visibleRows is how many data rows fit in a terminal height rows tall
// (at least 3), never more than there are rows.
func visibleRows(rows, height int) int {
	return max(min(rows, height-browseChrome), 3)
}

// headerHeight is the header line plus the rule under it. bubbles/table
// counts it as part of the height it's given.
const headerHeight = 2

// tableHeight is the height to give bubbles/table so that visibleRows data
// rows show below the header.
func tableHeight(rows, height int) int {
	return visibleRows(rows, height) + headerHeight
}

// tableWidth is the inner width of a bubbles/table with cols: each column plus
// one cell of padding on either side.
func tableWidth(cols []table.Column) int {
	w := 0
	for _, c := range cols {
		w += c.Width + 2
	}

	return w
}

func (m browseModel) Init() tea.Cmd { return nil }

func (m browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.done = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.table.SetHeight(tableHeight(len(m.table.Rows()), msg.Height-m.extra))
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m browseModel) View() tea.View {
	out := m.header + baseStyle.Render(m.table.View()) + "\n"
	// The help line goes away on quit, so what's left behind is just the table.
	if !m.done {
		out += "  " + m.table.HelpView() + dimStyle.Render(" • q quit") + "\n"
	}

	return tea.NewView(out)
}

// maxColumnWidth caps any one column so a long page title or URL can't
// stretch the table across a wide terminal.
const maxColumnWidth = 50

// fitColumns sizes the viewer's columns with columnWidths.
func fitColumns(t Table, width int) []table.Column {
	widths := columnWidths(t, width)
	cols := make([]table.Column, len(t.Headers))
	for i, h := range t.Headers {
		cols[i] = table.Column{Title: h, Width: widths[i]}
	}

	return cols
}

// columnWidths sizes each column to its widest cell (at most maxColumnWidth),
// then, when width is positive, shrinks the widest columns until the table
// fits in width. Each column also takes one cell of padding on both sides.
func columnWidths(t Table, width int) []int {
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = lipgloss.Width(h)
	}
	for _, r := range t.Rows {
		for i := 0; i < len(r) && i < len(widths); i++ {
			widths[i] = max(widths[i], lipgloss.Width(cleanCell(r[i])))
		}
	}
	for i := range widths {
		widths[i] = min(widths[i], maxColumnWidth)
	}

	if width <= 0 {
		return widths
	}

	avail := width - 2*len(widths)
	for sum(widths) > avail {
		widest := 0
		for i := range widths {
			if widths[i] > widths[widest] {
				widest = i
			}
		}
		if widths[widest] <= 4 {
			break
		}
		widths[widest]--
	}

	return widths
}

// fitCells returns a copy of t with every header and cell cut to its column
// width from columnWidths, so static output fits in t.Width. The outer border
// takes two columns.
func fitCells(t Table) Table {
	width := 0
	if t.Width > 0 {
		width = t.Width - 2
	}
	widths := columnWidths(t, width)

	fit := t
	fit.Headers = make([]string, len(t.Headers))
	for i, h := range t.Headers {
		fit.Headers[i] = ansi.Truncate(h, widths[i], "…")
	}
	fit.Rows = make([][]string, len(t.Rows))
	for r, row := range t.Rows {
		fit.Rows[r] = make([]string, len(row))
		for i, c := range row {
			c = cleanCell(c)
			if i < len(widths) {
				c = ansi.Truncate(c, widths[i], "…")
			}
			fit.Rows[r][i] = c
		}
	}

	return fit
}

// cleanCell flattens line breaks and tabs so one cell stays on one line.
func cleanCell(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func sum(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}

	return total
}

// WriteStatic writes RenderStatic(t, color) to w.
func WriteStatic(w io.Writer, t Table, color bool) error {
	_, err := io.WriteString(w, RenderStatic(t, color))
	return err
}
