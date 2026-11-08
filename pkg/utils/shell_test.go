package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendTokenToShellFileAppendsToExistingFile(t *testing.T) {
	assert := require.New(t)

	t.Setenv("SHELL", "/bin/zsh")
	t.Setenv("HOME", t.TempDir())

	rcPath := filepath.Join(os.Getenv("HOME"), ".zshenv")
	assert.NoError(os.WriteFile(rcPath, []byte("#!/bin/zsh\n"), 0o644))

	assert.NoError(AppendTokenToShellFile("TOKEN1"))

	contents, err := os.ReadFile(rcPath)
	assert.NoError(err)
	assert.Contains(string(contents), "#!/bin/zsh\n")
	assert.Contains(string(contents), `export LINKWARDEN_TOKEN="TOKEN1"`)
}

func TestAppendTokenToShellFileCreatesMissingFile(t *testing.T) {
	assert := require.New(t)

	t.Setenv("SHELL", "/bin/bash")
	t.Setenv("HOME", t.TempDir())

	assert.NoError(AppendTokenToShellFile("TOKEN2"))

	contents, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".bashrc"))
	assert.NoError(err)
	assert.Contains(string(contents), `export LINKWARDEN_TOKEN="TOKEN2"`)
}

func TestAppendTokenToShellFileCreatesFishConfigDir(t *testing.T) {
	assert := require.New(t)

	t.Setenv("SHELL", "/usr/local/bin/fish")
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", "")

	assert.NoError(AppendTokenToShellFile("TOKEN3"))

	contents, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish"))
	assert.NoError(err)
	assert.Contains(string(contents), `export LINKWARDEN_TOKEN="TOKEN3"`)
}

func TestAppendTokenToShellFileFishHonorsXDGConfigHome(t *testing.T) {
	assert := require.New(t)

	configHome := t.TempDir()
	t.Setenv("SHELL", "/usr/local/bin/fish")
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", configHome)

	assert.NoError(AppendTokenToShellFile("TOKEN5"))

	contents, err := os.ReadFile(filepath.Join(configHome, "fish", "config.fish"))
	assert.NoError(err)
	assert.Contains(string(contents), `export LINKWARDEN_TOKEN="TOKEN5"`)
}

func TestAppendTokenToShellFileUnrecognizedShellErrors(t *testing.T) {
	assert := require.New(t)

	t.Setenv("SHELL", "/bin/tcsh")
	t.Setenv("HOME", t.TempDir())

	assert.Error(AppendTokenToShellFile("TOKEN4"))
}

func TestShellConfigPathResolvesNonStandardShellPaths(t *testing.T) {
	assert := require.New(t)

	t.Setenv("SHELL", "/opt/homebrew/bin/zsh")
	t.Setenv("HOME", "/home/example")

	assert.Equal("/home/example/.zshenv", ShellConfigPath())
}
