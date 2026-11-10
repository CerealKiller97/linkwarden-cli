## linkwarden links bulk-delete

Delete multiple links

```bash
linkwarden links bulk-delete [flags]
```

### Options

```
  -h, --help[=true|false]   help for bulk-delete
      --ids ints            Link IDs. Comma-separated or repeated (--ids 1,2 --ids 3).
  -y, --yes[=true|false]    Skip the confirmation prompt (required when there's no interactive terminal).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links
