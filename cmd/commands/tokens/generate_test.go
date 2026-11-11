package tokens

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestStringToInt(t *testing.T) {
	assert := require.New(t)

	assert.Equal(SevenDays, stringToInt("7 days"))
	assert.Equal(ThirtyDays, stringToInt("30 days"))
	assert.Equal(SixtyDays, stringToInt("60 days"))
	assert.Equal(NinetyDays, stringToInt("90 days"))
	assert.Equal(NoExpiry, stringToInt("no expiration"))
	assert.Equal(NoExpiry, stringToInt("unrecognized value"))
}

func TestExecuteFailsWithoutInteractiveTerminal(t *testing.T) {
	assert := require.New(t)

	// generate has no flags at all (name/expiration are always prompted), so with
	// no interactive terminal available the first prompt fails immediately. This
	// at least exercises that the command wires up without panicking.
	utils.SetVars("test-token", "http://unused.invalid")

	cmd := Generate(nil)
	assert.Error(cmd.Execute())
}
