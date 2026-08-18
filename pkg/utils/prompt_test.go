package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveID_UsesFlagValueWithoutFetching(t *testing.T) {
	assert := require.New(t)

	called := false
	id, err := ResolveID("42", "--id", "Pick", func() ([]PickerItem, error) {
		called = true

		return nil, nil
	})

	assert.NoError(err)
	assert.Equal(42, id)
	assert.False(called)
}

func TestResolveID_InvalidFlagValue(t *testing.T) {
	assert := require.New(t)

	_, err := ResolveID("not-a-number", "--id", "Pick", func() ([]PickerItem, error) {
		return nil, nil
	})

	assert.Error(err)
}

// TestResolveID_NonInteractive_ReturnsErrorWithoutFetching covers the scripting/CI
// case: with no flag value and no interactive terminal to prompt from (true for any
// automated test run), ResolveID must fail fast with an actionable error naming the
// flag, rather than calling fetch (a network round-trip) or blocking on a picker
// nobody can answer.
func TestResolveID_NonInteractive_ReturnsErrorWithoutFetching(t *testing.T) {
	assert := require.New(t)

	called := false
	_, err := ResolveID("", "--id", "Pick", func() ([]PickerItem, error) {
		called = true

		return []PickerItem{{Label: "a", Value: 1}}, nil
	})

	assert.False(called, "fetch must not run when there's no terminal to prompt on")
	assert.Error(err)
	assert.Contains(err.Error(), "--id")
}
