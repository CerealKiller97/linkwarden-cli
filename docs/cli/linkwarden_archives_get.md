## linkwarden archives get

Retrieve a link's archive file

```bash
linkwarden archives get [flags]
```

### Examples

```bash
# Check whether a PDF snapshot exists (prints its type and size)
linkwarden archives get --link-id 412 --format pdf

# Save it to a file
linkwarden archives get --link-id 412 --format pdf --output effective-go.pdf

# Pick the link from a list and save its screenshot (PNG is the default)
linkwarden archives get --output screenshot.png
```

### Options

```
      --format string          Archive format. One of: png|jpeg|pdf|html|readable. (default "png")
  -h, --help[=true|false]      help for get
      --link-id int            Link ID.
      --output string          Path to save the archive file to. Without it, only content-type/size are printed.
      --preview[=true|false]   Retrieve the preview instead of the main archive file.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden archives](linkwarden_archives.md)	 - Manage link archive files
