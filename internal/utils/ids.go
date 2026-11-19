package utils

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// ErrInvalidID is returned when an ID flag is set to zero or a negative number.
	ErrInvalidID = errors.New("IDs must be positive")

	// ErrNoIDs is returned when a list-of-IDs flag was given no IDs.
	ErrNoIDs = errors.New("needs at least one ID")
)

// AddIDsFlag adds an --ids flag that takes comma-separated IDs, repeated
// flags, or both: --ids 1,2 --ids 3.
func AddIDsFlag(cmd *cobra.Command, usage string) {
	cmd.Flags().IntSlice("ids", nil, usage+" Comma-separated or repeated (--ids 1,2 --ids 3).")
}

// GetIDs returns the IDs passed to cmd's --ids flag. It fails when there are
// none, or when one isn't positive.
func GetIDs(cmd *cobra.Command) ([]int, error) {
	ids, err := cmd.Flags().GetIntSlice("ids")
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("--ids %w", ErrNoIDs)
	}

	for _, id := range ids {
		if id <= 0 {
			return nil, fmt.Errorf("%w: --ids %d", ErrInvalidID, id)
		}
	}

	return ids, nil
}
