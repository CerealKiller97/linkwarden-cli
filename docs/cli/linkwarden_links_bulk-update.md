## linkwarden links bulk-update

Apply a collection/tag change to multiple links

```
linkwarden links bulk-update [flags]
```

### Options

```
      --collection-id int                   Move links to this collection ID.
  -h, --help[=true|false]                   help for bulk-update
      --ids string                          Comma-separated list of link IDs to update.
      --remove-previous-tags[=true|false]   Remove tags not in --tag before applying.
      --tag strings                         Tag name to apply (repeatable).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links

