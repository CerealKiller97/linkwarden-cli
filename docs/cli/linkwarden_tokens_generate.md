## linkwarden tokens generate

Create an API token

### Synopsis

Create an API token for another machine or script.

Prompts for a token name and an expiry (7, 30, 60 or 90 days, or never), then
prints the new token. On Linux and macOS it offers to add the token to your
shell config. Needs an interactive terminal.

```
linkwarden tokens generate [flags]
```

### Options

```
  -h, --help[=true|false]   help for generate
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden tokens](linkwarden_tokens.md)	 - Manage tokens

