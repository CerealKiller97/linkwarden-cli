package tui

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInput_TypeAndConfirm(t *testing.T) {
	assert := require.New(t)

	m := press(newInputModel("Name", nil), append(typeText("ci"), keyEnter)...).(inputModel)

	assert.True(m.done)
	assert.Equal("ci", m.input.Value())
}

func TestInput_ValidationBlocksUntilFixed(t *testing.T) {
	assert := require.New(t)

	validate := func(s string) error {
		if s == "" {
			return errors.New("name is required")
		}
		return nil
	}

	m := press(newInputModel("Name", validate), keyEnter).(inputModel)
	assert.False(m.done)
	assert.Contains(m.View().Content, "name is required")

	m = press(m, append(typeText("x"), keyEnter)...).(inputModel)
	assert.True(m.done)
	assert.Equal("x", m.input.Value())
}

func TestInput_Esc_Aborts(t *testing.T) {
	m := press(newInputModel("Name", nil), keyEsc).(inputModel)
	require.True(t, m.aborted)
}
