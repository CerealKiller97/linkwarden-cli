## linkwarden links import

Import links from a JSON file

### Synopsis

Import links from a JSON file into one collection.

The file is a JSON array of objects with a "url" and a "title":

  [
    { "url": "https://go.dev/doc/effective_go", "title": "Effective Go" },
    { "url": "https://12factor.net", "title": "The Twelve-Factor App" }
  ]

Links already in the collection are skipped, as are URLs repeated within the
file. A summary of how many were created, skipped and failed is printed at the
end, and the command fails if any link could not be created.

Without --collection-id, a picker asks for the target collection.

```bash
linkwarden links import [flags]
```

### Examples

```bash
# Pick the target collection from a list
linkwarden links import --file reading-list.json

# Import into collection 3
linkwarden links import --file reading-list.json --collection-id 3

# Create 8 links at a time instead of 4
linkwarden links import --file reading-list.json --collection-id 3 --concurrency 8
```

### Options

```
      --collection-id int   Collection ID to import into.
      --concurrency int     Number of links to create in parallel. (default 4)
      --file string         JSON file of links to import.
  -h, --help[=true|false]   help for import
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden links](linkwarden_links.md)	 - Manage links
