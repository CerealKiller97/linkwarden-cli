package utils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// idCmd returns a command with an int --id flag, parsed from args.
func idCmd(t *testing.T, args ...string) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Int("id", 0, "ID.")
	require.NoError(t, cmd.Flags().Parse(args))

	return cmd
}

func TestResolveIDUsesFlagValueWithoutFetching(t *testing.T) {
	assert := require.New(t)

	called := false
	id, err := ResolveID(nil, idCmd(t, "--id", "42"), "id", "Pick", func() ([]PickerItem, error) {
		called = true

		return nil, nil
	})

	assert.NoError(err)
	assert.Equal(42, id)
	assert.False(called)
}

func TestResolveIDRejectsNonNumericFlagValue(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().Int("id", 0, "ID.")

	require.Error(t, cmd.Flags().Parse([]string{"--id", "not-a-number"}))
}

func TestResolveIDRejectsNonPositiveFlagValue(t *testing.T) {
	for _, value := range []string{"0", "-3"} {
		t.Run(value, func(t *testing.T) {
			assert := require.New(t)

			called := false
			_, err := ResolveID(nil, idCmd(t, "--id="+value), "id", "Pick", func() ([]PickerItem, error) {
				called = true

				return nil, nil
			})

			assert.ErrorIs(err, ErrInvalidID)
			assert.Contains(err.Error(), "--id")
			assert.False(called)
		})
	}
}

// TestResolveIDNonInteractiveReturnsErrorWithoutFetching covers the scripting/CI
// case: with no flag value and no interactive terminal to prompt from (true for any
// automated test run), ResolveID must fail fast with an actionable error naming the
// flag, rather than calling fetch (a network round-trip) or blocking on a picker
// nobody can answer.
func TestResolveIDNonInteractiveReturnsErrorWithoutFetching(t *testing.T) {
	assert := require.New(t)

	called := false
	_, err := ResolveID(nil, idCmd(t), "id", "Pick", func() ([]PickerItem, error) {
		called = true

		return []PickerItem{{Label: "a", Value: 1}}, nil
	})

	assert.False(called, "fetch must not run when there's no terminal to prompt on")
	assert.Error(err)
	assert.Contains(err.Error(), "--id")
}

func TestIsInteractiveFor(t *testing.T) {
	assert := require.New(t)

	assert.False(isInteractiveFor(true, true), "--non-interactive always wins, even on a real terminal")
	assert.False(isInteractiveFor(true, false), "--non-interactive wins, not a terminal either")
	assert.True(isInteractiveFor(false, true), "no override, real terminal -> interactive")
	assert.False(isInteractiveFor(false, false), "no override, not a terminal -> not interactive")
}

// TestResolveIDNonInteractiveOverrideReturnsErrorWithoutFetching exercises the
// --non-interactive root flag's plumbing end-to-end through SetNonInteractive, the
// same path cmd/root.go's PersistentPreRunE uses.
func TestResolveIDNonInteractiveOverrideReturnsErrorWithoutFetching(t *testing.T) {
	assert := require.New(t)

	SetNonInteractive(true)
	t.Cleanup(func() { SetNonInteractive(false) })

	called := false
	_, err := ResolveID(nil, idCmd(t), "id", "Pick", func() ([]PickerItem, error) {
		called = true

		return []PickerItem{{Label: "a", Value: 1}}, nil
	})

	assert.False(called)
	assert.Error(err)
	assert.Contains(err.Error(), "--id")
}
