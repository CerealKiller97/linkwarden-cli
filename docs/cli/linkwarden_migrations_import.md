## linkwarden migrations import

Import Linkwarden data from a data source

```
linkwarden migrations import [flags]
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

