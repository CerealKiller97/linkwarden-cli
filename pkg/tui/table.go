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

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
)

// RenderStatic returns t as a rounded-border text table. With color false the
// output has no ANSI escapes at all, which keeps piped output clean.
func RenderStatic(t Table, color bool) string {
	tbl := lgtable.New().
		Border(lipgloss.RoundedBorder()).
		Headers(t.Headers...).
		Rows(t.Rows...).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == lgtable.HeaderRow && color {
				return headerStyle.Padding(0, 1)
			}
			return cellStyle
		})

	title := t.Title
	if color {
		title = titleStyle.Render(title)
	}

	out := title + "\n" + tbl.String() + "\n"
	if !color {
		out = ansi.Strip(out)
	}

	return out
}

// TableChromeHeight is how many terminal rows RenderStatic adds around the data
// rows: title, top border, header, header separator and bottom border.
const TableChromeHeight = 5

// BrowseTable opens t in a scrollable bubbles/table sized to width×height (the
// terminal's size). It returns when the user quits with q, Esc or Ctrl+C.
func BrowseTable(t Table, width, height int, o Options) error {
	m := newBrowseModel(t, width, height)
	_, err := run(m, o)
	if errors.Is(err, ErrAborted) {
		return nil
	}

	return err
}

type browseModel struct {
	title string
	total int
	table table.Model
}

// browseChrome is the rows the viewer itself needs: title line, header,
// header border, and the footer.
const browseChrome = 4

func newBrowseModel(t Table, width, height int) browseModel {
	cols := fitColumns(t, width)
	rows := make([]table.Row, len(t.Rows))
	for i, r := range t.Rows {
		rows[i] = table.Row(r)
	}

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		Bold(true).
		Foreground(lipgloss.Color("14")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)
	styles.Selected = styles.Selected.
		Bold(true).
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("14"))

	tbl := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(max(height-browseChrome, 3)),
		// Without a width the viewport is zero columns wide and shows no rows.
		table.WithWidth(width),
		table.WithStyles(styles),
	)

	return browseModel{title: t.Title, total: len(t.Rows), table: tbl}
}

func (m browseModel) Init() tea.Cmd { return nil }

func (m browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(max(msg.Height-browseChrome, 3))
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m browseModel) View() tea.View {
	footer := dimStyle.Render(fmt.Sprintf("%d/%d · ↑/↓ j/k scroll · pgup/pgdn page · g/G top/bottom · q quit",
		m.table.Cursor()+1, m.total))

	v := tea.NewView(titleStyle.Render(m.title) + "\n" + m.table.View() + "\n" + footer)
	v.AltScreen = true

	return v
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
