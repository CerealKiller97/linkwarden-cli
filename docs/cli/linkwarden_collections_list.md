## linkwarden collections list

List collections

```bash
linkwarden collections list [flags]
```

### Examples

```bash
# List your collections
linkwarden collections list

# ID and name of each, tab-separated
linkwarden collections list --json | jq -r '.[] | "\(.id)\t\(.name)"'
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

* [linkwarden collections](linkwarden_collections.md)	 - Manage collections
