// Command gendocs writes the command reference in docs/cli/ from the Cobra command
// tree, so it always matches the flags. Run it after changing a command or flag:
//
//	go run ./scripts/gendocs
//
// The pre-commit hook in .githooks/ and CI fail when docs/cli/ is out of date.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CerealKiller97/linkwarden-cli/cmd"
	"github.com/spf13/cobra"
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

	return writeTree(root)
}

// writeTree writes one page per command, like doc.GenMarkdownTree, with the
// usage line and examples highlighted as bash.
func writeTree(c *cobra.Command) error {
	for _, child := range c.Commands() {
		if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := writeTree(child); err != nil {
			return err
		}
	}

	var buf bytes.Buffer
	if err := doc.GenMarkdownCustom(c, &buf, func(s string) string { return s }); err != nil {
		return err
	}

	name := strings.ReplaceAll(c.CommandPath(), " ", "_") + ".md"
	page := highlightShell(buf.String(), c)

	return os.WriteFile(filepath.Join(outDir, name), []byte(page), 0o600)
}

// highlightShell tags the usage and Examples code blocks as bash. Cobra writes
// every block as a bare ``` fence; the Options blocks stay that way, since they
// list flags rather than commands.
func highlightShell(page string, c *cobra.Command) string {
	if c.Runnable() {
		page = strings.Replace(page, "```\n"+c.UseLine()+"\n```", "```bash\n"+c.UseLine()+"\n```", 1)
	}
	if c.Example != "" {
		page = strings.Replace(page,
			"### Examples\n\n```\n"+c.Example+"\n```",
			"### Examples\n\n```bash\n"+dedent(c.Example)+"\n```", 1)
	}

	// Cobra ends each page with a blank line; keep just the final newline.
	return strings.TrimRight(page, "\n") + "\n"
}

// dedent removes the indentation every non-blank line shares. Examples are
// indented two spaces to line up in --help, which a code block doesn't need.
func dedent(s string) string {
	lines := strings.Split(s, "\n")

	indent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		n := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent == -1 || n < indent {
			indent = n
		}
	}
	if indent <= 0 {
		return s
	}

	for i, line := range lines {
		if len(line) >= indent {
			lines[i] = line[indent:]
		} else {
			lines[i] = strings.TrimLeft(line, " \t")
		}
	}

	return strings.Join(lines, "\n")
}
