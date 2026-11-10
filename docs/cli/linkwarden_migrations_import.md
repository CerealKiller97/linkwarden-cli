## linkwarden migrations import

Import Linkwarden data from a data source

```bash
linkwarden migrations import [flags]
```

### Examples

```bash
# Restore a Linkwarden export
linkwarden migrations import --file linkwarden-2026-10-06.json --format linkwarden

# Browser bookmarks exported as HTML (Chrome, Firefox, Safari…)
linkwarden migrations import --file bookmarks.html --format html

# Pocket export
linkwarden migrations import --file pocket.json --format pocket
```

### Options

```
      --file string         Path to the file containing the data to import.
      --format string       Import source format. One of: linkwarden|html|wallabag|omnivore|pocket.
  -h, --help[=true|false]   help for import
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden migrations](linkwarden_migrations.md)	 - Manage migration (Import / Export)
