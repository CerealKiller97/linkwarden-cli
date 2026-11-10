package choice

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type color int

var colors = []Option[color]{
	{Name: "red", Value: 7, Description: "Warm"},
	{Name: "blue", Value: 9, Description: "Cool"},
}

// run executes a command with a --color choice flag and returns what Get saw.
func run(t *testing.T, def string, args ...string) (color, bool, error) {
	t.Helper()

	var (
		got color
		set bool
	)
	cmd := &cobra.Command{
		Use: "paint",
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error
			got, set, err = Get[color](cmd, "color")

			return err
		},
	}
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	Flag(cmd, "color", def, "Paint color.", colors)
	cmd.SetArgs(args)

	return got, set, cmd.Execute()
}

func TestGetMapsNameToValue(t *testing.T) {
	assert := require.New(t)

	got, set, err := run(t, "", "--color", "blue")
	assert.NoError(err)
	assert.True(set)
	assert.Equal(color(9), got)
}

func TestNamesAreCaseInsensitive(t *testing.T) {
	assert := require.New(t)

	got, _, err := run(t, "", "--color", "RED")
	assert.NoError(err)
	assert.Equal(color(7), got)
}

func TestUnsetWithoutDefaultReportsNotSet(t *testing.T) {
	assert := require.New(t)

	_, set, err := run(t, "")
	assert.NoError(err)
	assert.False(set)
}

func TestUnsetUsesDefault(t *testing.T) {
	assert := require.New(t)

	got, set, err := run(t, "blue")
	assert.NoError(err)
	assert.True(set)
	assert.Equal(color(9), got)
}

func TestRejectsUnknownNamesAndNumbers(t *testing.T) {
	for _, arg := range []string{"green", "7", ""} {
		_, _, err := run(t, "", "--color", arg)
		require.ErrorContains(t, err, "must be one of red|blue", "arg %q", arg)
	}
}

func TestHelpListsNames(t *testing.T) {
	cmd := &cobra.Command{Use: "paint", Run: func(*cobra.Command, []string) {}}
	Flag(cmd, "color", "", "Paint color.", colors)

	require.Contains(t, cmd.Flags().Lookup("color").Usage, "Paint color. One of: red|blue.")
}

func TestShellCompletionOffersNamesWithDescriptions(t *testing.T) {
	assert := require.New(t)

	root := &cobra.Command{Use: "root"}
	paint := &cobra.Command{Use: "paint", Run: func(*cobra.Command, []string) {}}
	Flag(paint, "color", "", "Paint color.", colors)
	root.AddCommand(paint)

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{cobra.ShellCompRequestCmd, "paint", "--color", ""})
	assert.NoError(root.Execute())

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	assert.Equal([]string{"red\tWarm", "blue\tCool", ":4"}, lines[:3])
}
