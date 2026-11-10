## linkwarden links archive

Trigger archiving of a link

```bash
linkwarden links archive [flags]
```

### Examples

```bash
# Archive link 412 again (fresh screenshot, PDF and HTML)
linkwarden links archive --id 412

# Pick the link from a list
linkwarden links archive
```

### Options

```
  -h, --help[=true|false]   help for archive
      --id int              Link ID.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links
