package utils

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func idsCmd(t *testing.T, args ...string) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{Use: "test"}
	AddIDsFlag(cmd, "Link IDs.")
	require.NoError(t, cmd.Flags().Parse(args))

	return cmd
}

func TestGetIDsCommaSeparated(t *testing.T) {
	ids, err := GetIDs(idsCmd(t, "--ids", "1,2,3"))

	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, ids)
}

func TestGetIDsRepeatedAndMixed(t *testing.T) {
	ids, err := GetIDs(idsCmd(t, "--ids", "1,2", "--ids", "3"))

	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, ids)
}

func TestGetIDsEmpty(t *testing.T) {
	_, err := GetIDs(idsCmd(t))

	require.ErrorIs(t, err, ErrNoIDs)
}

func TestGetIDsRejectsNonPositive(t *testing.T) {
	_, err := GetIDs(idsCmd(t, "--ids", "4,0"))

	require.ErrorIs(t, err, ErrInvalidID)
}

func TestGetIDsRejectsNonNumeric(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	AddIDsFlag(cmd, "Link IDs.")

	require.Error(t, cmd.Flags().Parse([]string{"--ids", "1,abc"}))
}
