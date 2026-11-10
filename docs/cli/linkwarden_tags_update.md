## linkwarden tags update

Rename a tag

```bash
linkwarden tags update [flags]
```

### Examples

```bash
# Rename tag 14
linkwarden tags update --id 14 --name kubernetes

# Pick the tag from a list and rename it
linkwarden tags update --name kubernetes
```

### Options

```
  -h, --help[=true|false]   help for update
      --id int              Tag ID.
      --name string         New tag name.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden tags](linkwarden_tags.md)	 - Manage tags
