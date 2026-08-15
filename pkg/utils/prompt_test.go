package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveID_UsesFlagValueWithoutFetching(t *testing.T) {
	assert := require.New(t)

	called := false
	id, err := ResolveID("42", "Pick", func() ([]PickerItem, error) {
		called = true

		return nil, nil
	})

	assert.NoError(err)
	assert.Equal(42, id)
	assert.False(called)
}

func TestResolveID_InvalidFlagValue(t *testing.T) {
	assert := require.New(t)

	_, err := ResolveID("not-a-number", "Pick", func() ([]PickerItem, error) {
		return nil, nil
	})

	assert.Error(err)
}

func TestResolveID_PropagatesFetchError(t *testing.T) {
	assert := require.New(t)

	wantErr := errors.New("boom")
	_, err := ResolveID("", "Pick", func() ([]PickerItem, error) {
		return nil, wantErr
	})

	assert.ErrorIs(err, wantErr)
}

func TestResolveID_FetchesAndPromptsWhenFlagEmpty(t *testing.T) {
	assert := require.New(t)

	called := false
	_, err := ResolveID("", "Pick", func() ([]PickerItem, error) {
		called = true

		return []PickerItem{{Label: "a", Value: 1}}, nil
	})

	assert.True(called)
	// promptui needs an interactive terminal; in a non-interactive test run it
	// errors out immediately, which is expected here — this test only asserts that
	// omitting the flag triggers the picker path (via fetch being called) rather
	// than exercising an actual interactive selection.
	assert.Error(err)
}
