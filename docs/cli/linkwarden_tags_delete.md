## linkwarden tags delete

Delete a tag

```bash
linkwarden tags delete [flags]
```

### Examples

```bash
# Pick the tag from a list, then confirm
linkwarden tags delete

# Delete tag 14 without asking
linkwarden tags delete --id 14 --yes
```

### Options

```
  -h, --help[=true|false]   help for delete
      --id int              Tag ID.
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
