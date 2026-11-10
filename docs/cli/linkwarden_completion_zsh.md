## linkwarden completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(linkwarden completion zsh)

To load completions for every new session, execute once:

#### Linux:

	linkwarden completion zsh > "${fpath[1]}/_linkwarden"

#### macOS:

	linkwarden completion zsh > $(brew --prefix)/share/zsh/site-functions/_linkwarden

You will need to start a new shell for this setup to take effect.


```
linkwarden completion zsh [flags]
```

### Options

```
  -h, --help[=true|false]              help for zsh
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

