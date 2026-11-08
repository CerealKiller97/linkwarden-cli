package tui

import (
	"fmt"
	"io"
	"sync"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
)

// Progress is a bubbles/progress bar running as its own Bubble Tea program.
// It's safe to call Increment from many goroutines.
//
// The bar doesn't read the keyboard and doesn't install a signal handler, so
// Ctrl+C still interrupts the command the way it always has.
type Progress struct {
	program *tea.Program
	done    chan struct{}
	once    sync.Once
	// final is the model the program ended with; set before done is closed.
	final tea.Model
}

// StartProgress draws a bar for total units of work under label on out
// (normally stderr, so stdout stays clean for results).
func StartProgress(label string, total int, out io.Writer) *Progress {
	return startProgress(label, total, out)
}

func startProgress(label string, total int, out io.Writer, extra ...tea.ProgramOption) *Progress {
	p := &Progress{done: make(chan struct{})}
	opts := append([]tea.ProgramOption{
		tea.WithInput(nil),
		tea.WithOutput(out),
		tea.WithoutSignalHandler(),
	}, extra...)
	p.program = tea.NewProgram(newProgressModel(label, total), opts...)

	go func() {
		defer close(p.done)
		p.final, _ = p.program.Run()
	}()

	return p
}

// Increment records one finished unit of work.
func (p *Progress) Increment() { p.program.Send(progressIncrMsg{}) }

// Stop finishes the bar and waits for its last frame to be drawn.
func (p *Progress) Stop() {
	p.once.Do(func() {
		p.program.Send(progressStopMsg{})
		<-p.done
	})
}

type (
	progressIncrMsg struct{}
	progressStopMsg struct{}
)

type progressModel struct {
	label   string
	total   int
	current int
	bar     progress.Model
	stopped bool
}

func newProgressModel(label string, total int) progressModel {
	return progressModel{
		label: label,
		total: total,
		bar:   progress.New(progress.WithColors(Ocean, Sky), progress.WithWidth(40)),
	}
}

func (m progressModel) Init() tea.Cmd { return nil }

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressIncrMsg:
		m.current = min(m.current+1, m.total)
		return m, m.bar.SetPercent(m.percent())
	case progressStopMsg:
		m.stopped = true
		return m, tea.Quit
	case progress.FrameMsg:
		var cmd tea.Cmd
		m.bar, cmd = m.bar.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m progressModel) percent() float64 {
	if m.total == 0 {
		return 1
	}

	return float64(m.current) / float64(m.total)
}

func (m progressModel) View() tea.View {
	// Once stopped, draw the bar at its real value rather than mid-animation.
	bar := m.bar.View()
	if m.stopped {
		bar = m.bar.ViewAs(m.percent())
	}

	return tea.NewView(fmt.Sprintf("%s %s %d/%d\n", m.label, bar, m.current, m.total))
}
