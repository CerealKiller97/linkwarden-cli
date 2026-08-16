package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// getShell returns the path (relative to $HOME) of the file the current shell reads
// environment variables from, or "" if the shell isn't recognized (e.g. on Windows).
func getShell() string {
	shell := filepath.Base(os.Getenv("SHELL"))

	switch shell {
	case "zsh":
		// .zshenv is sourced for every zsh invocation (login, interactive, and
		// non-interactive/scripts), unlike .zshrc which only loads interactively.
		return ".zshenv"
	case "bash":
		// bash has no equivalent of .zshenv read unconditionally by every shell.
		// .bashrc is the best available choice: most distros' .bash_profile sources
		// it, so interactive login and non-login shells both pick it up.
		return ".bashrc"
	case "fish":
		// fish has a single config file read for both login and interactive shells.
		return filepath.Join(".config", "fish", "config.fish")
	default:
		return "" // Windows, or an unrecognized shell. Not supported.
	}
}

// ShellConfigPath returns the absolute path to the current shell's config file, or
// "" if the shell isn't recognized.
func ShellConfigPath() string {
	shellFile := getShell()
	if shellFile == "" {
		return ""
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, shellFile)
}

// AppendTokenToShellFile appends an `export LINKWARDEN_TOKEN=...` line to the
// current shell's config file, creating the file (and any parent directory, e.g.
// ~/.config/fish) if it doesn't already exist.
func AppendTokenToShellFile(token string) error {
	path := ShellConfigPath()
	if path == "" {
		return errors.New("unrecognized or unsupported shell")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	// 0600: this file may now contain a secret token, so keep it readable only by
	// the owner. If the file already existed with looser permissions, those are
	// left as-is (O_CREATE only applies the mode on creation).
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, "\n# Managed by Linkwarden CLI\nexport LINKWARDEN_TOKEN=\"%s\"\n", token); err != nil {
		return err
	}

	return nil
}
