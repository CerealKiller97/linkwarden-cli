package utils

import (
	"os"
	"path/filepath"
)

// appDirName is the subdirectory the CLI uses inside each XDG base directory.
const appDirName = "linkwarden"

// ConfigHome returns $XDG_CONFIG_HOME, falling back to ~/.config.
func ConfigHome() (string, error) {
	return xdgBaseDir("XDG_CONFIG_HOME", ".config")
}

// StateHome returns $XDG_STATE_HOME, falling back to ~/.local/state.
func StateHome() (string, error) {
	return xdgBaseDir("XDG_STATE_HOME", ".local", "state")
}

// StateFile returns the path of the named file in the CLI's state directory,
// $XDG_STATE_HOME/linkwarden/<name>. It does not create the directory.
func StateFile(name string) (string, error) {
	dir, err := StateHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, appDirName, name), nil
}

// xdgBaseDir resolves an XDG base directory the same way on every OS, Windows
// included: the env var's value when it is an absolute path, otherwise (unset,
// empty or relative, all of which the XDG spec says to ignore) the default under
// the user's home directory ($HOME on Unix, %USERPROFILE% on Windows).
func xdgBaseDir(envVar string, defaultUnderHome ...string) (string, error) {
	if dir := os.Getenv(envVar); dir != "" && filepath.IsAbs(dir) {
		return dir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(append([]string{home}, defaultUnderHome...)...), nil
}
