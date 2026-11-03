package tokens

import (
	"testing"

	"github.com/CerealKiller97/linkwarden-cli/internal/e2etest"
	"github.com/stretchr/testify/require"
)

func TestRevoke_FailsWithoutInteractiveTerminal(t *testing.T) {
	assert := require.New(t)

	e2etest.Setup(t)

	// Revoke always picks the token interactively. A test run has no terminal, so it
	// fetches the tokens and then fails at the picker without revoking anything.
	cmd := Revoke()
	assert.Error(cmd.Execute())

	_, err := e2etest.Setup(t).Collections.List(t.Context())
	assert.NoError(err, "the test's own token must still work")
}

func TestRevoke_PropagatesListError(t *testing.T) {
	assert := require.New(t)

	e2etest.UseInvalidToken(t)

	cmd := Revoke()
	assert.Error(cmd.Execute())
}
