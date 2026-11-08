package utils

import (
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

// InstanceStateFile returns the path of the state file for one Linkwarden
// instance, $XDG_STATE_HOME/linkwarden/<name>/<instance>, so state such as a
// collection ID from one instance is never used against another. <instance> is
// instanceURL's host, port and path (see instanceKey). It does not create the
// directory.
func InstanceStateFile(name, instanceURL string) (string, error) {
	dir, err := StateFile(name)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, instanceKey(instanceURL)), nil
}

// unsafeFileNameChars matches anything that isn't safe in a file name on every
// OS (":" isn't allowed on Windows, "/" separates directories).
var unsafeFileNameChars = regexp.MustCompile(`[^a-z0-9.-]+`)

// instanceKey turns an instance URL into a file name: its host, port and path,
// lowercased, with every other character run replaced by "_". The scheme is
// dropped, since http and https on one host are the same instance:
// "https://cloud.linkwarden.app" is "cloud.linkwarden.app", and
// "http://localhost:3100" is "localhost_3100".
func instanceKey(instanceURL string) string {
	key := instanceURL
	if u, err := url.Parse(instanceURL); err == nil && u.Host != "" {
		key = u.Host + u.Path
	}

	key = unsafeFileNameChars.ReplaceAllString(strings.ToLower(key), "_")
	key = strings.Trim(key, "_.")
	if key == "" {
		return "default"
	}

	return key
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
