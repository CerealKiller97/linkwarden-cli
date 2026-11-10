## linkwarden collections create

Create a collection

### Synopsis

Create a collection

```bash
linkwarden collections create [flags]
```

### Examples

```bash
# Create a collection
linkwarden collections create --name "Reading List" --color "#f59e0b"

# Nest a collection inside collection 9
linkwarden collections create --name Kubernetes --parent-id 9 --icon cube

# Create one and keep its ID for a script
COLLECTION=$(linkwarden collections create --name Inbox -q)
```

### Options

```
      --color string         Collection color. (default "#00bba7")
      --description string   Collection description. (default "Managed by: Linkwarden CLI")
  -h, --help[=true|false]    help for create
      --icon string          Collection icon.
      --icon-weight string   Collection icon weight.
      --name string          Collection name.
      --parent-id int        Parent collection ID. Defaults to the Unorganized collection.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden collections](linkwarden_collections.md)	 - Manage collections
