## linkwarden collections get

Get a collection by ID

```bash
linkwarden collections get [flags]
```

### Examples

```bash
# Pick the collection from a list
linkwarden collections get

# Get collection 7
linkwarden collections get --id 7

# Print only its name
linkwarden collections get --id 7 --json | jq -r .name
```

### Options

```
  -h, --help[=true|false]   help for get
      --id int              Collection ID.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden collections](linkwarden_collections.md)	 - Manage collections
