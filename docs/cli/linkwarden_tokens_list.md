## linkwarden tokens list

List all available tokens

```bash
linkwarden tokens list [flags]
```

### Examples

```bash
# List your API tokens
linkwarden tokens list

# Same, as JSON
linkwarden tokens list --json
```

### Options

```
  -h, --help[=true|false]   help for list
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden tokens](linkwarden_tokens.md)	 - Manage tokens
