package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func getShell() string {
	shell := os.Getenv("SHELL")

	switch shell {
	case "/bin/zsh":
		return ".zshrc" // or .zshenv
	case "/bin/bash":
		return ".bashrc" // or profile
	case "/bin/fish":
		return ".fishrc"
	default:
		return "" // Windows. Not Supported :)
	}
}
func AppendTokenToShellFile(token string) error {
	shellFile := getShell()

	if shellFile == "" {
		return errors.New("shell file not exist") // Windows
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join(dir, shellFile), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		if err := f.Close(); err != nil {
		}
	}(f)

	template := `
	# Managed by Linkwarden CLI 
	export LINKWARDEN_TOKEN="%s"
	`
	formattedTemplate := fmt.Sprintf(template, token)

	_, err = f.WriteString(formattedTemplate)
	if err != nil {
		return err
	}

	return nil
}

func Reload() error {
	// Source: Trust me bro
	// source is a Bash shell built-in command that executes the content of the file passed as an argument in the current shell.
	// It has a synonym in . (period).
	cmd := exec.Command("source", "~/"+getShell())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
