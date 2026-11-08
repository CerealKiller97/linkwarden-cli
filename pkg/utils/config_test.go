package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigFileIsNamespacedUnderLinkwarden(t *testing.T) {
	assert := require.New(t)

	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	got, err := ConfigFile()
	assert.NoError(err)
	assert.Equal(filepath.Join(dir, "linkwarden", "config.json"), got)
}

func TestLoadConfigMissingFileReturnsEmpty(t *testing.T) {
	assert := require.New(t)

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := LoadConfig()
	assert.NoError(err)
	assert.Equal(Config{}, cfg)
}

func TestSaveConfigRoundTripsWithOwnerOnlyPermissions(t *testing.T) {
	assert := require.New(t)

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	want := Config{URL: "https://links.example.com", Token: "secret"}
	path, err := SaveConfig(want)
	assert.NoError(err)

	got, err := LoadConfig()
	assert.NoError(err)
	assert.Equal(want, got)

	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		assert.NoError(err)
		assert.Equal(os.FileMode(0o600), info.Mode().Perm())
	}
}

func TestLoadConfigInvalidJSONErrors(t *testing.T) {
	assert := require.New(t)

	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	assert.NoError(os.MkdirAll(filepath.Join(dir, "linkwarden"), 0o700))
	assert.NoError(os.WriteFile(filepath.Join(dir, "linkwarden", "config.json"), []byte("{"), 0o600))

	_, err := LoadConfig()
	assert.Error(err)
}
