## linkwarden dashboard update-v2

Update dashboard layout (v2)

```bash
linkwarden dashboard update-v2 [flags]
```

### Examples

```bash
# layout.json: [{"type":"STATS","enabled":true,"order":0},{"type":"RECENT_LINKS","enabled":true,"order":1}]
linkwarden dashboard update-v2 --file layout.json
```

### Options

```
      --file string         Path to a JSON file containing the layout, e.g. [{"type":"STATS","enabled":true,"order":0}].
  -h, --help[=true|false]   help for update-v2
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden dashboard](linkwarden_dashboard.md)	 - View and configure the dashboard
