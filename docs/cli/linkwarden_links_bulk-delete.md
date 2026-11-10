## linkwarden links bulk-delete

Delete multiple links

```bash
linkwarden links bulk-delete [flags]
```

### Examples

```bash
# Delete two links (asks first)
linkwarden links bulk-delete --ids 301,302

# Same, without asking (scripts, cron)
linkwarden links bulk-delete --ids 301 --ids 302 --yes

# Delete every link that matches a search
linkwarden links list --search "old draft" -q | paste -sd, - | xargs -I{} linkwarden links bulk-delete --ids {} --yes
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
