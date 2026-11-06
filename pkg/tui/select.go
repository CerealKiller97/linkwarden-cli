package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// defaultPageSize is how many options are visible at once.
const defaultPageSize = 10

// Select shows a filterable list of options under label and returns the index
// (into options) of the one the user picks. Typing narrows the list, ↑/↓ (or
// Ctrl+P/Ctrl+N) move, Enter confirms, Esc/Ctrl+C aborts with ErrAborted.
func Select(label string, options []string) (int, error) {
	return SelectWith(label, options, Options{})
}

// SelectWith is Select with an explicit terminal.
func SelectWith(label string, options []string, o Options) (int, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("%s: nothing to choose from", strings.TrimSuffix(label, ":"))
	}

	final, err := run(newSelectModel(label, options), o)
	if err != nil {
		return 0, err
	}

	m := final.(selectModel)
	if m.aborted {
		return 0, ErrAborted
	}

	return m.choice, nil
}

type selectModel struct {
	label    string
	options  []string
	pageSize int

	filter  string
	matches []int // indexes into options matching filter
	cursor  int   // position within matches
	offset  int   // first visible position within matches
	choice  int
	done    bool
	aborted bool
}

func newSelectModel(label string, options []string) selectModel {
	m := selectModel{label: label, options: options, pageSize: defaultPageSize, choice: -1}
	m.applyFilter()

	return m
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	switch key.String() {
	case "ctrl+c", "esc":
		m.aborted = true
		return m, tea.Quit
	case "enter":
		if len(m.matches) == 0 {
			return m, nil
		}
		m.choice = m.matches[m.cursor]
		m.done = true
		return m, tea.Quit
	case "up", "ctrl+p":
		m.move(-1)
	case "down", "ctrl+n":
		m.move(1)
	case "pgup":
		m.move(-m.pageSize)
	case "pgdown":
		m.move(m.pageSize)
	case "home":
		m.move(-len(m.matches))
	case "end":
		m.move(len(m.matches))
	case "backspace":
		if r := []rune(m.filter); len(r) > 0 {
			m.filter = string(r[:len(r)-1])
			m.applyFilter()
		}
	case "ctrl+u":
		m.filter = ""
		m.applyFilter()
	default:
		if key.Text != "" {
			m.filter += key.Text
			m.applyFilter()
		}
	}

	return m, nil
}

func (m *selectModel) move(delta int) {
	if len(m.matches) == 0 {
		return
	}

	m.cursor = min(max(m.cursor+delta, 0), len(m.matches)-1)

	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.pageSize {
		m.offset = m.cursor - m.pageSize + 1
	}
}

func (m *selectModel) applyFilter() {
	needle := strings.ToLower(m.filter)
	m.matches = make([]int, 0, len(m.options))
	for i, opt := range m.options {
		if needle == "" || strings.Contains(strings.ToLower(opt), needle) {
			m.matches = append(m.matches, i)
		}
	}
	m.cursor, m.offset = 0, 0
}

func (m selectModel) View() tea.View {
	var b strings.Builder

	if m.done {
		fmt.Fprintf(&b, "%s %s %s\n", doneStyle.Render("✔"), labelStyle.Render(m.label), m.options[m.choice])
		return tea.NewView(b.String())
	}
	if m.aborted {
		return tea.NewView("")
	}

	b.WriteString(labelStyle.Render(m.label))
	if m.filter != "" {
		b.WriteString(" " + m.filter)
	} else {
		b.WriteString(" " + dimStyle.Render("type to filter"))
	}
	b.WriteString("\n")

	if len(m.matches) == 0 {
		b.WriteString(dimStyle.Render("  no matches") + "\n")
	}

	end := min(m.offset+m.pageSize, len(m.matches))
	for pos := m.offset; pos < end; pos++ {
		opt := m.options[m.matches[pos]]
		if pos == m.cursor {
			b.WriteString(cursorStyle.Render("› ") + selectedStyle.Render(opt) + "\n")
		} else {
			b.WriteString("  " + opt + "\n")
		}
	}

	footer := "↑/↓ move · enter select · esc cancel"
	if len(m.matches) > m.pageSize {
		footer = fmt.Sprintf("%d/%d · %s", m.cursor+1, len(m.matches), footer)
	}
	b.WriteString(dimStyle.Render(footer) + "\n")

	return tea.NewView(b.String())
}
