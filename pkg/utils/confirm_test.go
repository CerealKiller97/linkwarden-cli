package utils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func newConfirmCommand(t *testing.T, args ...string) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{Use: "test"}
	AddYesFlag(cmd)
	require.NoError(t, cmd.ParseFlags(args))

	return cmd
}

func TestConfirmWithYesSkipsPrompt(t *testing.T) {
	assert := require.New(t)

	for _, flag := range []string{"--yes", "-y"} {
		called := false
		ok, err := ConfirmWith(newConfirmCommand(t, flag), func() (string, error) {
			called = true

			return "Delete?", nil
		})

		assert.NoError(err, flag)
		assert.True(ok, flag)
		assert.False(called, "%s must not build the question", flag)
	}
}

// TestConfirmWithoutTerminalRequiresYes covers scripts and CI: test runs have no
// terminal on stdin, so without --yes nothing may be deleted.
func TestConfirmWithoutTerminalRequiresYes(t *testing.T) {
	assert := require.New(t)

	called := false
	ok, err := ConfirmWith(newConfirmCommand(t), func() (string, error) {
		called = true

		return "Delete?", nil
	})

	assert.ErrorIs(err, ErrConfirmationRequired)
	assert.Contains(err.Error(), "--yes")
	assert.False(ok)
	assert.False(called, "no question is built when there's no one to ask")
}

func TestJoinIDsAndPlural(t *testing.T) {
	assert := require.New(t)

	assert.Equal("301, 302, 303", JoinIDs([]int{301, 302, 303}))
	assert.Equal("7", JoinIDs([]int{7}))
	assert.Equal("1 link", Plural(1, "link"))
	assert.Equal("3 links", Plural(3, "link"))
	assert.Equal("0 tags", Plural(0, "tag"))
}

func TestConfirmNonInteractiveOverrideRequiresYes(t *testing.T) {
	assert := require.New(t)

	SetNonInteractive(true)
	t.Cleanup(func() { SetNonInteractive(false) })

	ok, err := Confirm(newConfirmCommand(t), "Delete?")

	assert.ErrorIs(err, ErrConfirmationRequired)
	assert.False(ok)
}
