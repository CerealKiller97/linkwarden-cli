## linkwarden archives upload-standalone

Upload a local file as a new link's archive

```bash
linkwarden archives upload-standalone [flags]
```

### Examples

```bash
# Save a page you downloaded yourself as a new link
linkwarden archives upload-standalone --url https://go.dev/doc/effective_go --file effective-go.html --format html

# Same with a PDF
linkwarden archives upload-standalone --url https://12factor.net --file twelve-factor.pdf --format pdf
```

### Options

```
      --file string            Local file to upload.
      --format string          Archive format. One of: png|jpeg|pdf|html.
  -h, --help[=true|false]      help for upload-standalone
      --preview[=true|false]   Store the uploaded file as the generated preview instead of the main archive.
      --url string             URL to associate with the new link.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden archives](linkwarden_archives.md)	 - Manage link archive files
