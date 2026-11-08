package tui

import (
	"io"
	"sync"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// spinnerDelay is how long a call can run before the spinner appears, so a fast
// response doesn't flash a frame of it on screen.
const spinnerDelay = 150 * time.Millisecond

// Spinner is a bubbles/spinner running as its own Bubble Tea program while
// something without a known length runs, such as an API request.
//
// Like Progress, it doesn't read the keyboard and doesn't install a signal
// handler, so Ctrl+C still interrupts the command the way it always has.
type Spinner struct {
	program *tea.Program
	done    chan struct{}
	once    sync.Once
}

// StartSpinner draws a spinner next to label on out (normally stderr, so
// stdout stays clean for results). It clears itself when stopped.
func StartSpinner(label string, out io.Writer) *Spinner {
	return startSpinner(label, spinnerDelay, out)
}

func startSpinner(label string, delay time.Duration, out io.Writer, extra ...tea.ProgramOption) *Spinner {
	s := &Spinner{done: make(chan struct{})}
	opts := append([]tea.ProgramOption{
		tea.WithInput(nil),
		tea.WithOutput(out),
		tea.WithoutSignalHandler(),
	}, extra...)
	s.program = tea.NewProgram(newSpinnerModel(label, delay), opts...)

	go func() {
		defer close(s.done)
		_, _ = s.program.Run()
	}()

	return s
}

// Stop removes the spinner and waits until it's gone from the screen, so
// whatever is printed next starts on a clean line.
func (s *Spinner) Stop() {
	s.once.Do(func() {
		s.program.Send(spinnerStopMsg{})
		<-s.done
	})
}

type (
	spinnerShowMsg struct{}
	spinnerStopMsg struct{}
)

type spinnerModel struct {
	label   string
	delay   time.Duration
	spinner spinner.Model
	shown   bool
	stopped bool
}

func newSpinnerModel(label string, delay time.Duration) spinnerModel {
	return spinnerModel{
		label: label,
		delay: delay,
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(Sky)),
		),
	}
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.Tick(m.delay, func(time.Time) tea.Msg {
		return spinnerShowMsg{}
	}))
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinnerShowMsg:
		m.shown = true
		return m, nil
	case spinnerStopMsg:
		m.stopped = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m spinnerModel) View() tea.View {
	// An empty last frame erases the spinner rather than leaving it behind.
	if !m.shown || m.stopped {
		return tea.NewView("")
	}

	return tea.NewView(m.spinner.View() + " " + dimStyle.Render(m.label))
}
