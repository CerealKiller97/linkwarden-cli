package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// yesFlag is the flag destructive commands take to skip their confirmation prompt.
const yesFlag = "yes"

// ErrConfirmationRequired is returned by Confirm when there's no interactive
// terminal to ask on and --yes wasn't passed.
var ErrConfirmationRequired = errors.New(
	"this can't be undone and needs confirmation, but there's no interactive terminal: pass --yes (-y) to go ahead")

// AddYesFlag registers --yes / -y on a destructive command. Confirm reads it.
func AddYesFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP(yesFlag, "y", false,
		"Skip the confirmation prompt (required when there's no interactive terminal).")
}

// Confirm asks question as a No/Yes prompt before a destructive action and
// reports whether to go ahead. See ConfirmWith.
func Confirm(cmd *cobra.Command, question string) (bool, error) {
	return ConfirmWith(cmd, func() (string, error) { return question, nil })
}

// ConfirmWith is Confirm with a question that is only built when the prompt is
// actually shown, so commands can fetch details (a name, a link count) for it
// without paying for that request under --yes or in scripts.
//
// With --yes set it returns true without asking. Without an interactive
// terminal (or with --non-interactive) it returns ErrConfirmationRequired
// rather than deleting anything unasked. Otherwise "No" is the first, default
// answer, so a stray Enter cancels.
func ConfirmWith(cmd *cobra.Command, question func() (string, error)) (bool, error) {
	yes, err := cmd.Flags().GetBool(yesFlag)
	if err != nil {
		return false, err
	}

	if yes {
		return true, nil
	}

	if !isInteractive() {
		return false, ErrConfirmationRequired
	}

	q, err := question()
	if err != nil {
		return false, err
	}

	_, answer, err := Select(q, []string{"No", "Yes"})
	if err != nil {
		return false, fmt.Errorf("confirm: %w", err)
	}

	return answer == "Yes", nil
}

// JoinIDs formats ids for a confirmation question, e.g. "301, 302, 303".
func JoinIDs(ids []int) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.Itoa(id)
	}

	return strings.Join(parts, ", ")
}

// Plural returns "1 link" or "3 links".
func Plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", noun)
	}

	return fmt.Sprintf("%d %ss", n, noun)
}
