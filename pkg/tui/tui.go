// Package tui holds the CLI's interactive terminal prompts, built on Bubble Tea.
//
// The prompts here don't decide whether a terminal is available; callers (see
// pkg/utils) are expected to check that before running one.
package tui

import (
	"errors"
	"io"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ErrAborted is returned when the user cancels a prompt with Esc or Ctrl+C.
var ErrAborted = errors.New("prompt aborted")

var (
	labelStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	cursorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	dimStyle      = lipgloss.NewStyle().Faint(true)
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	doneStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

// Options let tests (or unusual callers) swap the terminal Bubble Tea talks to.
type Options struct {
	Input  io.Reader
	Output io.Writer
}

func (o Options) programOptions() []tea.ProgramOption {
	var opts []tea.ProgramOption
	if o.Input != nil {
		opts = append(opts, tea.WithInput(o.Input))
	}
	if o.Output != nil {
		opts = append(opts, tea.WithOutput(o.Output))
	}

	return opts
}

// run executes m and maps an interrupted program to ErrAborted.
func run(m tea.Model, o Options) (tea.Model, error) {
	final, err := tea.NewProgram(m, o.programOptions()...).Run()
	if errors.Is(err, tea.ErrInterrupted) {
		return final, ErrAborted
	}

	return final, err
}
