## linkwarden migrations export

Export Linkwarden data to a JSON file

```bash
linkwarden migrations export [flags]
```

### Options

```
  -h, --help[=true|false]   help for export
      --output string       Path to write the export to. Defaults to ./linkwarden-export-<timestamp>.json.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden migrations](linkwarden_migrations.md)	 - Manage migration (Import / Export)
