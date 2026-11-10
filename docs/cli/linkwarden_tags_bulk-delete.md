## linkwarden tags bulk-delete

Delete multiple tags

```bash
linkwarden tags bulk-delete [flags]
```

### Examples

```bash
# Delete three tags (asks first)
linkwarden tags bulk-delete --ids 14,15,22

# Same, without asking
linkwarden tags bulk-delete --ids 14 --ids 15 --ids 22 --yes
```

### Options

```
  -h, --help[=true|false]   help for bulk-delete
      --ids ints            Tag IDs. Comma-separated or repeated (--ids 1,2 --ids 3).
  -y, --yes[=true|false]    Skip the confirmation prompt (required when there's no interactive terminal).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden tags](linkwarden_tags.md)	 - Manage tags
