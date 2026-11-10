## linkwarden search

Search links

```bash
linkwarden search [flags]
```

### Examples

```bash
# Search your links
linkwarden search --query "kubernetes operators"

# Oldest matches first
linkwarden search --query golang --sort oldest

# IDs only, one per line
linkwarden search --query golang -q
```

### Options

```
      --cursor int          Pagination cursor.
  -h, --help[=true|false]   help for search
      --query string        Search query string.
      --sort string         Sort order (the server's default when omitted). One of: newest|oldest|name-asc|name-desc|description-asc|description-desc.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden](linkwarden.md)	 - Linkwarden is a CLI for the Linkwarden platform
