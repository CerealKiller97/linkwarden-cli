package utils

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// configFileName is the CLI's config file inside $XDG_CONFIG_HOME/linkwarden.
const configFileName = "config.json"

// Config is what `linkwarden init` saves: the instance URL and the access token
// used against it. The LINKWARDEN_URL / LINKWARDEN_TOKEN env vars override it.
type Config struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// ConfigFile returns the path of the CLI's config file,
// $XDG_CONFIG_HOME/linkwarden/config.json. It does not create the directory.
func ConfigFile() (string, error) {
	dir, err := ConfigHome()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, appDirName, configFileName), nil
}

// LoadConfig reads the config file. A missing file is not an error: it returns
// an empty Config.
func LoadConfig() (Config, error) {
	path, err := ConfigFile()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path) //nolint:gosec // path is the CLI's own config file
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// SaveConfig writes cfg to the config file, creating its directory if needed.
// The file holds a secret token, so it is readable only by the owner.
func SaveConfig(cfg Config) (string, error) {
	path, err := ConfigFile()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return "", err
	}

	// WriteFile only applies the mode on creation; tighten an existing file too.
	if err := os.Chmod(path, 0o600); err != nil {
		return "", err
	}

	return path, nil
}
