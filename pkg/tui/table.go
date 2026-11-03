package tui

import (
	"errors"
	"fmt"
	"io"

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
}

// The look follows Bubble Tea's table example: a thin gray outer border, a
// gray rule under the header, no column separators, and a purple highlight on
// the selected row.
var (
	borderColor = lipgloss.Color("240")
	baseStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(borderColor)
	titleStyle  = lipgloss.NewStyle().Bold(true)
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
)

// RenderStatic returns t as plain text in the same style as the viewer. With
// color false the output has no ANSI escapes at all, so piped output is clean.
func RenderStatic(t Table, color bool) string {
	tbl := lgtable.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(borderColor)).
		BorderColumn(false).
		BorderHeader(true).
		Headers(t.Headers...).
		Rows(t.Rows...).
		StyleFunc(func(_, _ int) lipgloss.Style { return cellStyle })

	out := titleStyle.Render(t.Title) + "\n" + tbl.String() + "\n"
	if !color {
		out = ansi.Strip(out)
	}

	return out
}

// TableChromeHeight is how many terminal rows RenderStatic adds around the data
// rows: title, top border, header, header rule and bottom border.
const TableChromeHeight = 5

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
	title string
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
	rows := make([]table.Row, len(t.Rows))
	for i, r := range t.Rows {
		rows[i] = table.Row(r)
	}

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		BorderBottom(true).
		Bold(false)
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	tbl := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		// Styles before height: WithHeight measures the (styled) header.
		table.WithStyles(styles),
		table.WithHeight(tableHeight(len(rows), height)),
		// Without a width the viewport is zero columns wide and shows no rows.
		table.WithWidth(tableWidth(cols)),
	)

	return browseModel{title: fmt.Sprintf("%s (%d)", t.Title, len(t.Rows)), table: tbl}
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
		m.table.SetHeight(tableHeight(len(m.table.Rows()), msg.Height))
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m browseModel) View() tea.View {
	out := titleStyle.Render(m.title) + "\n" + baseStyle.Render(m.table.View()) + "\n"
	// The help line goes away on quit, so what's left behind is just the table.
	if !m.done {
		out += "  " + m.table.HelpView() + dimStyle.Render(" • q quit") + "\n"
	}

	return tea.NewView(out)
}

// fitColumns sizes each column to its widest cell, then shrinks the widest
// columns until the table fits in width.
func fitColumns(t Table, width int) []table.Column {
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = lipgloss.Width(h)
	}
	for _, r := range t.Rows {
		for i := 0; i < len(r) && i < len(widths); i++ {
			widths[i] = max(widths[i], lipgloss.Width(r[i]))
		}
	}

	// Each column gets one cell of padding on both sides from the default styles.
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

	cols := make([]table.Column, len(t.Headers))
	for i, h := range t.Headers {
		cols[i] = table.Column{Title: h, Width: widths[i]}
	}

	return cols
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
