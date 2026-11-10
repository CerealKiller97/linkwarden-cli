## linkwarden config get

Get public runtime configuration

```bash
linkwarden config get [flags]
```

### Examples

```bash
# Show your instance's public configuration
linkwarden config get

# Same, as JSON
linkwarden config get --json
```

### Options

```
  -h, --help[=true|false]   help for get
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden config](linkwarden_config.md)	 - View instance configuration
