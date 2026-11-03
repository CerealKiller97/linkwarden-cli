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

	m := final.(inputModel)
	if m.aborted {
		return "", ErrAborted
	}

	return m.input.Value(), nil
}

type inputModel struct {
	label    string
	validate func(string) error
	input    textinput.Model
	err      error
	done     bool
	aborted  bool
}

func newInputModel(label string, validate func(string) error) inputModel {
	ti := textinput.New()
	ti.Prompt = ""
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
				if err := m.validate(m.input.Value()); err != nil {
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
		b.WriteString(doneStyle.Render("✔") + " " + labelStyle.Render(m.label) + " " + m.input.Value() + "\n")
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
