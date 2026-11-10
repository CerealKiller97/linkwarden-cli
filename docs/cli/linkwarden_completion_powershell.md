## linkwarden completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	linkwarden completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```bash
linkwarden completion powershell [flags]
```

### Options

```
  -h, --help[=true|false]              help for powershell
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
