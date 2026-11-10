## linkwarden collections update

Update a collection

```bash
linkwarden collections update [flags]
```

### Examples

```bash
# Rename collection 7
linkwarden collections update --id 7 --name "Weekend Reads"

# Make it public, then private again
linkwarden collections update --id 7 --public
linkwarden collections update --id 7 --public=false

# Pick a collection from a list and move it inside collection 9
linkwarden collections update --parent-id 9
```

### Options

```
      --color string          New collection color.
      --description string    New collection description.
  -h, --help[=true|false]     help for update
      --icon string           New collection icon.
      --icon-weight string    New collection icon weight.
      --id int                Collection ID.
      --name string           New collection name.
      --parent-id int         New parent collection ID.
      --public[=true|false]   Make the collection public (--public=false makes it private).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden collections](linkwarden_collections.md)	 - Manage collections
