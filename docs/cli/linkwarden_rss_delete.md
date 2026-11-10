## linkwarden rss delete

Delete an RSS subscription

```bash
linkwarden rss delete [flags]
```

### Examples

```bash
# Pick the subscription from a list, then confirm
linkwarden rss delete

# Delete subscription 5 without asking
linkwarden rss delete --id 5 --yes
```

### Options

```
  -h, --help[=true|false]   help for delete
      --id int              Subscription ID.
  -y, --yes[=true|false]    Skip the confirmation prompt (required when there's no interactive terminal).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden rss](linkwarden_rss.md)	 - Manage RSS subscriptions
