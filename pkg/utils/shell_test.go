package utils

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestName(t *testing.T) {
	// Arrange
	assert := require.New(t)

	t.Setenv("SHELL", "/bin/zsh")
	t.Setenv("HOME", t.TempDir())
	create, err := os.Create(path.Join(os.Getenv("HOME"), ".zshrc"))
	assert.NoError(err)
	writeString, err := create.WriteString("#!/bin/zsh\n")
	assert.NoError(err)
	assert.GreaterOrEqual(writeString, 1)
	err = create.Close()
	assert.NoError(err)

	// Act
	err = AppendTokenToShellFile("TOKEN1")

	file, err := os.ReadFile(path.Join(os.Getenv("HOME"), ".zshrc"))
	assert.NoError(err)
	assert.Contains(string(file), "export LINKWARDEN_TOKEN=\"TOKEN1\"")

	// Assert
	assert.NoError(err)
}
