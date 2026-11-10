## linkwarden links delete

Delete a link

```bash
linkwarden links delete [flags]
```

### Examples

```bash
# Pick the link from a list, then confirm
linkwarden links delete

# Delete link 412 without asking (scripts, cron)
linkwarden links delete --id 412 --yes
```

### Options

```
  -h, --help[=true|false]   help for delete
      --id int              Link ID.
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
