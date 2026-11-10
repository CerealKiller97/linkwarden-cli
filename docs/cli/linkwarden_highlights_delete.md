## linkwarden highlights delete

Delete a highlight

```bash
linkwarden highlights delete [flags]
```

### Examples

```bash
# Pick a link, then one of its highlights
linkwarden highlights delete

# Pick from the highlights on link 412
linkwarden highlights delete --link-id 412

# Delete highlight 31 without asking
linkwarden highlights delete --id 31 --yes
```

### Options

```
  -h, --help[=true|false]   help for delete
      --id int              Highlight ID.
      --link-id int         Link ID (used to pick a highlight when --id is not set).
  -y, --yes[=true|false]    Skip the confirmation prompt (required when there's no interactive terminal).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden highlights](linkwarden_highlights.md)	 - Manage link highlights
