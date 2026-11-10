## linkwarden links list

List links

```bash
linkwarden links list [flags]
```

### Examples

```bash
# Your newest links
linkwarden links list

# Links in collection 3
linkwarden links list --collection-id 3

# Pinned links with tag 14
linkwarden links list --tag-id 14 --pinned-only

# Links that match a search
linkwarden links list --search kubernetes

# The URL of each link, one per line
linkwarden links list --json | jq -r '.[].url'
```

### Options

```
      --collection-id int          Filter by collection ID.
      --cursor int                 Pagination cursor.
  -h, --help[=true|false]          help for list
      --pinned-only[=true|false]   Only show pinned links.
      --search string              Search query string.
      --tag-id int                 Filter by tag ID.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links
