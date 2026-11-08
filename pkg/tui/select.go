package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// defaultPageSize is how many options are visible at once: one paginator page.
const defaultPageSize = 10

// SelectOption is one choice in a picker: a label, and optional dim detail
// shown after it, such as a link's shortened URL.
type SelectOption struct {
	Label  string
	Detail string
}

// LoadMore fetches the options after the ones already loaded. more reports
// whether there may be further options after these.
type LoadMore func() (options []SelectOption, more bool, err error)

// Select shows a filterable list of options under label and returns the index
// (into options) of the one the user picks. Typing narrows the list, ↑/↓ (or
// Ctrl+P/Ctrl+N) move, ←/→ (or PgUp/PgDn) change page, Enter confirms,
// Esc/Ctrl+C aborts with ErrAborted.
func Select(label string, options []string) (int, error) {
	return SelectWith(label, options, Options{})
}

// SelectWith is Select with an explicit terminal.
func SelectWith(label string, options []string, o Options) (int, error) {
	opts := make([]SelectOption, len(options))
	for i, opt := range options {
		opts[i] = SelectOption{Label: opt}
	}

	return SelectPaged(label, opts, nil, o)
}

// SelectPaged is Select over options loaded a batch at a time: first is the
// batch already loaded, and when the user reaches the last page of what's
// loaded, load (if non-nil) fetches the next batch in the background. Typing
// filters the options loaded so far, and keeps loading while the matches fit
// on one page, so a filter can find options that weren't loaded yet. It
// returns the picked option's index across every batch, in load order.
func SelectPaged(label string, first []SelectOption, load LoadMore, o Options) (int, error) {
	if len(first) == 0 {
		return 0, fmt.Errorf("%s: nothing to choose from", strings.TrimSuffix(label, ":"))
	}

	final, err := run(newPagedSelectModel(label, first, load), o)
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
	label   string
	options []SelectOption
	pages   paginator.Model
	width   int

	// load fetches the next batch; more is false once it has said there are
	// no more, or failed. loading is true while a fetch is in flight.
	load    LoadMore
	more    bool
	loading bool
	loadErr error

	filter  string
	matches []int // indexes into options matching filter
	cursor  int   // position within matches
	choice  int
	done    bool
	aborted bool
}

// loadedMsg carries the result of a LoadMore call back to Update.
type loadedMsg struct {
	options []SelectOption
	more    bool
	err     error
}

func newSelectModel(label string, options []string) selectModel {
	opts := make([]SelectOption, len(options))
	for i, opt := range options {
		opts[i] = SelectOption{Label: opt}
	}

	return newPagedSelectModel(label, opts, nil)
}

func newPagedSelectModel(label string, options []SelectOption, load LoadMore) selectModel {
	pages := paginator.New(paginator.WithPerPage(defaultPageSize))
	pages.ArabicFormat = "page %d/%d"

	m := selectModel{label: label, options: options, pages: pages, load: load, more: load != nil, choice: -1}
	m.applyFilter()

	return m
}

// Init starts loading right away when everything loaded so far fits on one
// page, so there's a second page to move to.
func (m selectModel) Init() tea.Cmd {
	return m.maybeLoad()
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case loadedMsg:
		m.loading = false
		if msg.err != nil {
			m.loadErr = msg.err
			m.more = false
			return m, nil
		}
		m.options = append(m.options, msg.options...)
		m.more = msg.more && len(msg.options) > 0
		m.refilter()
		return m, m.maybeLoad()
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m selectModel) handleKey(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
	case "left", "pgup":
		m.move(-m.pages.PerPage)
	case "right", "pgdown":
		m.move(m.pages.PerPage)
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
		if key.Text == "" {
			return m, nil
		}
		m.filter += key.Text
		m.applyFilter()
	}

	return m, m.maybeLoad()
}

// maybeLoad starts fetching the next batch when the user is on the last page
// of the matches and more may exist, unless a fetch is already running.
func (m *selectModel) maybeLoad() tea.Cmd {
	if m.load == nil || !m.more || m.loading || !m.pages.OnLastPage() {
		return nil
	}

	m.loading = true
	load := m.load

	return func() tea.Msg {
		opts, more, err := load()
		return loadedMsg{options: opts, more: more, err: err}
	}
}

func (m *selectModel) move(delta int) {
	if len(m.matches) == 0 {
		return
	}

	m.cursor = min(max(m.cursor+delta, 0), len(m.matches)-1)
	m.syncPage()
}

// syncPage points the paginator at the page holding the cursor.
func (m *selectModel) syncPage() {
	m.pages.TotalPages = max(1, (len(m.matches)+m.pages.PerPage-1)/m.pages.PerPage)
	m.pages.Page = m.cursor / m.pages.PerPage
}

// applyFilter recomputes the matches for a changed filter and goes back to
// the first one.
func (m *selectModel) applyFilter() {
	m.cursor = 0
	m.refilter()
}

// refilter recomputes the matches, keeping the cursor where it is (newly
// loaded options only ever add matches after the current ones).
func (m *selectModel) refilter() {
	needle := strings.ToLower(m.filter)
	m.matches = make([]int, 0, len(m.options))
	for i, opt := range m.options {
		if needle == "" || strings.Contains(strings.ToLower(opt.Label+" "+opt.Detail), needle) {
			m.matches = append(m.matches, i)
		}
	}
	m.cursor = min(m.cursor, max(len(m.matches)-1, 0))
	m.syncPage()
}

func (m selectModel) View() tea.View {
	var b strings.Builder

	if m.done {
		fmt.Fprintf(&b, "%s %s %s\n", doneStyle.Render("✔"), labelStyle.Render(m.label), m.options[m.choice].Label)
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

	if len(m.matches) == 0 && !m.loading {
		b.WriteString(dimStyle.Render("  no matches") + "\n")
	}

	start, end := m.pages.GetSliceBounds(len(m.matches))
	for pos := start; pos < end; pos++ {
		b.WriteString(m.row(m.options[m.matches[pos]], pos == m.cursor) + "\n")
	}

	b.WriteString(m.footer() + "\n")

	return tea.NewView(b.String())
}

// minLabelWidth is the narrowest a label is cut to before the detail is
// dropped instead, so a long URL never squeezes out the name.
const minLabelWidth = 20

// row renders one option: a cursor mark, the label and its dim detail, cut to
// the terminal width when it's known.
func (m selectModel) row(opt SelectOption, selected bool) string {
	label, detail := opt.Label, opt.Detail

	if m.width > 0 {
		avail := m.width - 2 // the cursor mark
		detailWidth := 0
		if detail != "" {
			detailWidth = 2 + lipgloss.Width(detail)
		}

		switch {
		case lipgloss.Width(label)+detailWidth <= avail:
		case avail-detailWidth >= minLabelWidth:
			label = ansi.Truncate(label, avail-detailWidth, "…")
		default:
			detail = ""
			label = ansi.Truncate(label, avail, "…")
		}
	}

	line := label
	if selected {
		line = selectedStyle.Render(label)
	}
	if detail != "" {
		line += "  " + dimStyle.Render(detail)
	}

	if selected {
		return cursorStyle.Render("› ") + line
	}

	return "  " + line
}

// footer is the page indicator, loading state and key help.
func (m selectModel) footer() string {
	var parts []string

	if m.pages.TotalPages > 1 || m.more {
		page := m.pages.View()
		if m.more {
			page += "+"
		}
		parts = append(parts, page)
	}

	switch {
	case m.loading:
		parts = append(parts, "loading more…")
	case m.loadErr != nil:
		parts = append(parts, errorStyle.Render("couldn't load more: "+m.loadErr.Error()))
	}

	help := "↑/↓ move · enter select · esc cancel"
	if m.pages.TotalPages > 1 || m.more {
		help = "↑/↓ move · ←/→ page · enter select · esc cancel"
	}
	parts = append(parts, help)

	return dimStyle.Render(strings.Join(parts, " · "))
}
