## linkwarden links delete-archive

Delete archived files for multiple links

```bash
linkwarden links delete-archive [flags]
```

### Examples

```bash
# Delete the archived files of two links (asks first)
linkwarden links delete-archive --ids 412,418

# Same, without asking
linkwarden links delete-archive --ids 412 --ids 418 --yes
```

### Options

```
  -h, --help[=true|false]   help for delete-archive
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
