## linkwarden links get

Get a link by ID

```bash
linkwarden links get [flags]
```

### Examples

```bash
# Pick the link from a list
linkwarden links get

# Get link 412
linkwarden links get --id 412

# Print only its URL
linkwarden links get --id 412 --json | jq -r .url
```

### Options

```
  -h, --help[=true|false]   help for get
      --id int              Link ID.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links
