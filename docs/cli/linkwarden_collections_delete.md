## linkwarden collections delete

Delete a collection

```bash
linkwarden collections delete [flags]
```

### Examples

```bash
# Pick the collection from a list, then confirm
linkwarden collections delete

# Delete collection 7 without asking (scripts, cron)
linkwarden collections delete --id 7 --yes
```

### Options

```
  -h, --help[=true|false]   help for delete
      --id int              Collection ID.
  -y, --yes[=true|false]    Skip the confirmation prompt (required when there's no interactive terminal).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden collections](linkwarden_collections.md)	 - Manage collections
