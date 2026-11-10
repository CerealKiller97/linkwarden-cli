## linkwarden highlights create

Create or update a highlight on a link

```bash
linkwarden highlights create [flags]
```

### Examples

```bash
# Highlight a passage on link 412
linkwarden highlights create --link-id 412 --text "Don't communicate by sharing memory." \
  --color yellow --start-offset 1840 --end-offset 1876

# Pick the link from a list and add a comment
linkwarden highlights create --text "Clear is better than clever." --color yellow \
  --start-offset 120 --end-offset 148 --comment "quote this"
```

### Options

```
      --color string        Highlight color.
      --comment string      Optional comment.
      --end-offset int      End offset of the highlighted text.
  -h, --help[=true|false]   help for create
      --link-id int         Link ID.
      --start-offset int    Start offset of the highlighted text.
      --text string         Highlighted text.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden highlights](linkwarden_highlights.md)	 - Manage link highlights
