## linkwarden links create

Save a link

### Synopsis

Save a single link.

Without --collection-id the link goes to your Unorganized collection, like in the
web app. Without --name, Linkwarden uses the page's title.

To save many links at once from a JSON file, use "links import".

```bash
linkwarden links create <url> [flags]
```

### Examples

```bash
# Save a link to Unorganized; its name comes from the page's title
linkwarden links create https://go.dev/doc/effective_go

# Name it and save it to collection 3
linkwarden links create https://12factor.net --name "The Twelve-Factor App" --collection-id 3

# Add tags
linkwarden links create https://k8spatterns.io --tag kubernetes --tag books
```

### Options

```
      --collection-id int    Collection ID to save to. Defaults to Unorganized.
      --description string   Link description.
  -h, --help[=true|false]    help for create
      --name string          Link name. Defaults to the page's title.
      --tag strings          Tag name to add (repeatable).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links
