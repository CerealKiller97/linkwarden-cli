// Package prompttest answers the CLI's interactive prompts from a script, so
// tests can drive commands that ask questions (init, tokens generate, any
// confirmation) without a terminal.
package prompttest

import (
	"strings"
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/tui"
	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
)

// Answer is the reply to one prompt. Label must appear in the prompt's label,
// so a test fails if the command asks something else than expected. Value is
// the option to pick for a Select, or the text to type otherwise; an empty
// Value for InputDefault accepts the default.
type Answer struct {
	Label string
	Value string

	index    int
	useIndex bool
}

// Pick answers a Select by position (0-based) rather than by text, for options
// that share a name.
func Pick(label string, index int) Answer {
	return Answer{Label: label, index: index, useIndex: true}
}

// Script is a Prompter that hands out its answers in order.
type Script struct {
	t       testing.TB
	answers []Answer
}

// New returns a Script that replies to the prompts the code under test shows,
// in order. Pass it to testutils.WithPrompter. A prompt with no answer left, or
// a label that doesn't match, fails the test and aborts the prompt. Answers
// left unused at the end fail the test too.
func New(t testing.TB, answers ...Answer) *Script {
	t.Helper()

	s := &Script{t: t, answers: answers}
	t.Cleanup(s.checkUsed)

	return s
}

// Use is New for tests that don't run with a testutils App yet: it makes the
// script answer every prompt until the test ends.
//
// Legacy, removed once every command reads the App: pass New's Script to testutils.WithPrompter instead.
func Use(t testing.TB, answers ...Answer) *Script {
	t.Helper()

	s := &Script{t: t, answers: answers}
	restore := utils.SetPrompter(s)
	t.Cleanup(func() {
		restore()
		s.checkUsed()
	})

	return s
}

func (s *Script) checkUsed() {
	if len(s.answers) > 0 {
		s.t.Errorf("prompttest: %d answer(s) never asked for, first %q", len(s.answers), s.answers[0].Label)
	}
}

func (s *Script) next(label string) (Answer, bool) {
	if len(s.answers) == 0 {
		s.t.Errorf("prompttest: unexpected prompt %q", label)
		return Answer{}, false
	}

	a := s.answers[0]
	s.answers = s.answers[1:]

	if !strings.Contains(label, a.Label) {
		s.t.Errorf("prompttest: got prompt %q, want one containing %q", label, a.Label)
		return Answer{}, false
	}

	return a, true
}

// Select picks the next answer's option, by text or, with Pick, by position.
func (s *Script) Select(label string, options []string) (int, error) {
	a, ok := s.next(label)
	if !ok {
		return 0, tui.ErrAborted
	}

	if a.useIndex {
		if a.index < 0 || a.index >= len(options) {
			s.t.Errorf("prompttest: %q has %d options, can't pick #%d", label, len(options), a.index)
			return 0, tui.ErrAborted
		}

		return a.index, nil
	}

	for i, o := range options {
		if o == a.Value {
			return i, nil
		}
	}

	s.t.Errorf("prompttest: %q isn't an option of %q (options: %q)", a.Value, label, options)

	return 0, tui.ErrAborted
}

// Input types the next answer.
func (s *Script) Input(label string, validate func(string) error) (string, error) {
	return s.text(label, "", validate)
}

// InputDefault types the next answer; an empty one accepts def.
func (s *Script) InputDefault(label string, def string, validate func(string) error) (string, error) {
	return s.text(label, def, validate)
}

// Password types the next answer.
func (s *Script) Password(label string, validate func(string) error) (string, error) {
	return s.text(label, "", validate)
}

// text types an answer the way the real prompt would take it: an empty answer
// takes def, and an answer the prompt's validation rejects fails the test,
// because the real prompt would keep asking.
func (s *Script) text(label string, def string, validate func(string) error) (string, error) {
	a, ok := s.next(label)
	if !ok {
		return "", tui.ErrAborted
	}

	v := a.Value
	if v == "" {
		v = def
	}

	if validate != nil {
		if err := validate(v); err != nil {
			s.t.Errorf("prompttest: answer %q to %q fails the prompt's validation: %v", v, label, err)
			return "", tui.ErrAborted
		}
	}

	return v, nil
}
