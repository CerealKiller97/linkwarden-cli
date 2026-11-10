// Command gendocs writes the command reference in docs/cli/ from the Cobra command
// tree, so it always matches the flags. Run it after changing a command or flag:
//
//	go run ./scripts/gendocs
//
// CI fails when docs/cli/ is out of date.
package main

import (
	"fmt"
	"os"

	"github.com/CerealKiller97/linkwarden-cli/cmd"
	"github.com/spf13/cobra/doc"
)

const outDir = "docs/cli"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gendocs:", err)
		os.Exit(1)
	}
}

func run() error {
	root := cmd.Root()
	// Cobra adds these lazily on Execute; add them now so they're documented too.
	root.InitDefaultHelpCmd()
	root.InitDefaultCompletionCmd()
	// No "Auto generated on <date>" footer, so regenerating unchanged commands is a no-op.
	root.DisableAutoGenTag = true

	// Start from an empty folder so removed commands don't leave stale pages behind.
	if err := os.RemoveAll(outDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		return err
	}

	return doc.GenMarkdownTree(root, outDir)
}
