## linkwarden completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	linkwarden completion fish | source

To load completions for every new session, execute once:

	linkwarden completion fish > ~/.config/fish/completions/linkwarden.fish

You will need to start a new shell for this setup to take effect.


```bash
linkwarden completion fish [flags]
```

### Options

```
  -h, --help[=true|false]              help for fish
      --no-descriptions[=true|false]   disable completion descriptions
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden completion](linkwarden_completion.md)	 - Generate the autocompletion script for the specified shell
