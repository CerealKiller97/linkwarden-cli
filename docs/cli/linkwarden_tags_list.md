## linkwarden tags list

List tags

```bash
linkwarden tags list [flags]
```

### Examples

```bash
# List your tags
linkwarden tags list

# Tags that match a search
linkwarden tags list --search go

# Tag IDs only, one per line
linkwarden tags list -q
```

### Options

```
      --cursor int          Pagination cursor.
  -h, --help[=true|false]   help for list
      --search string       Search query string.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden tags](linkwarden_tags.md)	 - Manage tags
