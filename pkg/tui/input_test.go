package tui

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
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

func TestInputDefault_PrefilledWithDefault(t *testing.T) {
	assert := require.New(t)

	m := newDefaultInputModel("URL", "https://cloud.linkwarden.app", nil)
	assert.Contains(m.View().Content, "ttps://cloud.linkwarden.app", "the whole default is on screen, not just its first letter")
	assert.Contains(m.View().Content, "ctrl+u clear")

	m = press(m, keyEnter).(inputModel)
	assert.True(m.done)
	assert.Equal("https://cloud.linkwarden.app", m.value())
}

func TestInputDefault_ClearedFieldAcceptsDefault(t *testing.T) {
	m := newDefaultInputModel("URL", "https://cloud.linkwarden.app", nil)

	m = press(m, keyCtrlU, keyEnter).(inputModel)
	require.Equal(t, "https://cloud.linkwarden.app", m.value())
}

func TestInputDefault_TypedValueWins(t *testing.T) {
	m := newDefaultInputModel("URL", "https://cloud.linkwarden.app", nil)

	m = press(m, append([]tea.KeyPressMsg{keyCtrlU}, append(typeText("http://localhost:3000"), keyEnter)...)...).(inputModel)
	require.Equal(t, "http://localhost:3000", m.value())
}

func TestInputDefault_EditsDefault(t *testing.T) {
	m := newDefaultInputModel("URL", "https://cloud.linkwarden.app", nil)

	m = press(m, append(typeText("/"), keyEnter)...).(inputModel)
	require.Equal(t, "https://cloud.linkwarden.app/", m.value(), "the cursor starts at the end of the default")
}

var keyCtrlU = tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl}

func TestPassword_NeverEchoesValue(t *testing.T) {
	assert := require.New(t)

	m := press(newPasswordModel("Password", nil), typeText("hunter2")...).(inputModel)
	assert.NotContains(m.View().Content, "hunter2")

	m = press(m, keyEnter).(inputModel)
	assert.True(m.done)
	assert.Equal("hunter2", m.value())
	assert.NotContains(m.View().Content, "hunter2")
}
