package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// setHome points os.UserHomeDir at dir on every OS ($HOME on Unix, %USERPROFILE%
// on Windows).
func setHome(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
}

func TestStateHome_UsesAbsoluteEnvVar(t *testing.T) {
	assert := require.New(t)

	stateDir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateDir)

	got, err := StateHome()
	assert.NoError(err)
	assert.Equal(stateDir, got)
}

func TestStateHome_FallsBackWhenUnsetEmptyOrRelative(t *testing.T) {
	for name, value := range map[string]string{
		"empty":    "",
		"relative": filepath.Join("relative", "state"),
	} {
		t.Run(name, func(t *testing.T) {
			assert := require.New(t)

			home := t.TempDir()
			setHome(t, home)
			t.Setenv("XDG_STATE_HOME", value)

			got, err := StateHome()
			assert.NoError(err)
			assert.Equal(filepath.Join(home, ".local", "state"), got)
		})
	}
}

func TestConfigHome_UsesEnvVarOrDefault(t *testing.T) {
	assert := require.New(t)

	home := t.TempDir()
	setHome(t, home)

	t.Setenv("XDG_CONFIG_HOME", "")
	got, err := ConfigHome()
	assert.NoError(err)
	assert.Equal(filepath.Join(home, ".config"), got)

	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	got, err = ConfigHome()
	assert.NoError(err)
	assert.Equal(configDir, got)
}

func TestStateFile_IsNamespacedUnderLinkwarden(t *testing.T) {
	assert := require.New(t)

	stateDir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateDir)

	got, err := StateFile("unorganized")
	assert.NoError(err)
	assert.Equal(filepath.Join(stateDir, "linkwarden", "unorganized"), got)
}

func TestInstanceStateFile_OneFilePerInstance(t *testing.T) {
	assert := require.New(t)

	stateDir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", stateDir)

	got, err := InstanceStateFile("unorganized", "https://cloud.linkwarden.app")
	assert.NoError(err)
	assert.Equal(filepath.Join(stateDir, "linkwarden", "unorganized", "cloud.linkwarden.app"), got)
}

func TestInstanceKey(t *testing.T) {
	cases := map[string]string{
		"https://cloud.linkwarden.app":          "cloud.linkwarden.app",
		"http://cloud.linkwarden.app":           "cloud.linkwarden.app",
		"http://localhost:3100":                 "localhost_3100",
		"http://127.0.0.1:3000":                 "127.0.0.1_3000",
		"https://Links.Example.com/linkwarden":  "links.example.com_linkwarden",
		"https://links.example.com/linkwarden/": "links.example.com_linkwarden",
		"http://[::1]:3000":                     "1_3000",
		"":                                      "default",
	}

	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			require.Equal(t, want, instanceKey(in))
		})
	}
}
