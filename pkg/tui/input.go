package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Input asks for a single line of text under label. validate, when non-nil, is
// checked on Enter; its error is shown inline and the prompt stays open until
// the value passes. Esc/Ctrl+C aborts with ErrAborted.
func Input(label string, validate func(string) error) (string, error) {
	return InputWith(label, validate, Options{})
}

// InputWith is Input with an explicit terminal.
func InputWith(label string, validate func(string) error, o Options) (string, error) {
	final, err := run(newInputModel(label, validate), o)
	if err != nil {
		return "", err
	}

	return inputResult(final)
}

// InputDefault is Input with def shown as a placeholder; pressing Enter on an
// empty field accepts def.
func InputDefault(label, def string, validate func(string) error) (string, error) {
	m := newInputModel(label, validate)
	m.def = def
	m.input.Placeholder = def

	final, err := run(m, Options{})
	if err != nil {
		return "", err
	}

	return inputResult(final)
}

// Password is Input with the typed characters masked, on screen and in the
// confirmation line left behind.
func Password(label string, validate func(string) error) (string, error) {
	final, err := run(newPasswordModel(label, validate), Options{})
	if err != nil {
		return "", err
	}

	return inputResult(final)
}

func inputResult(final tea.Model) (string, error) {
	m := final.(inputModel)
	if m.aborted {
		return "", ErrAborted
	}

	return m.value(), nil
}

func newPasswordModel(label string, validate func(string) error) inputModel {
	m := newInputModel(label, validate)
	m.masked = true
	m.input.EchoMode = textinput.EchoPassword
	m.input.EchoCharacter = '•'

	return m
}

// value is what the user entered, or def when they left the field empty.
func (m inputModel) value() string {
	if m.def != "" && strings.TrimSpace(m.input.Value()) == "" {
		return m.def
	}

	return m.input.Value()
}

type inputModel struct {
	label    string
	validate func(string) error
	def      string // accepted when Enter is pressed on an empty field
	masked   bool   // password input: never echo the value
	input    textinput.Model
	err      error
	done     bool
	aborted  bool
}

func newInputModel(label string, validate func(string) error) inputModel {
	ti := textinput.New()
	ti.Prompt = ""
	styles := textinput.DefaultDarkStyles()
	styles.Cursor.Color = Sky
	ti.SetStyles(styles)
	ti.Focus()

	return inputModel{label: label, validate: validate, input: ti}
}

func (m inputModel) Init() tea.Cmd { return textinput.Blink }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			m.aborted = true
			return m, tea.Quit
		case "enter":
			if m.validate != nil {
				if err := m.validate(m.value()); err != nil {
					m.err = err
					return m, nil
				}
			}
			m.done = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if _, ok := msg.(tea.KeyPressMsg); ok {
		m.err = nil
	}

	return m, cmd
}

func (m inputModel) View() tea.View {
	var b strings.Builder

	switch {
	case m.done:
		shown := m.value()
		if m.masked {
			shown = strings.Repeat("•", min(len([]rune(shown)), 8))
		}
		b.WriteString(doneStyle.Render("✔") + " " + labelStyle.Render(m.label) + " " + shown + "\n")
	case m.aborted:
	default:
		b.WriteString(labelStyle.Render(m.label) + " " + m.input.View() + "\n")
		if m.err != nil {
			b.WriteString(errorStyle.Render("✘ "+m.err.Error()) + "\n")
		} else {
			b.WriteString(dimStyle.Render("enter confirm · esc cancel") + "\n")
		}
	}

	return tea.NewView(b.String())
}
