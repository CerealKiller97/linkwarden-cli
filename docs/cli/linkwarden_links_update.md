## linkwarden links update

Update a link

```bash
linkwarden links update [flags]
```

### Examples

```bash
# Rename link 412
linkwarden links update --id 412 --name "Effective Go"

# Point it at a new URL
linkwarden links update --id 412 --url https://go.dev/doc/effective_go

# Pick the link from a list and change its description
linkwarden links update --description "The classic Go style guide"
```

### Options

```
      --color string         New link color.
      --description string   New link description.
  -h, --help[=true|false]    help for update
      --icon string          New link icon.
      --icon-weight string   New link icon weight.
      --id int               Link ID.
      --name string          New link name.
      --url string           New link URL.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links
