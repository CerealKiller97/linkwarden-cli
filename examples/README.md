# Scripting examples

Small scripts that use the global `--json` and `-q` / `--quiet` flags to drive `linkwarden` from the shell.
They need [`jq`](https://jqlang.org) (`brew install jq`) and a signed-in CLI (`linkwarden init`, or
`LINKWARDEN_TOKEN`).

The scripts run `linkwarden` from your `PATH`. To try them with a build of this repo instead:

```bash
go build -o bin/linkwarden ./cmd/linkwarden
export LINKWARDEN_BIN="$PWD/bin/linkwarden"
```

## Scripts

| Script | What it does |
| --- | --- |
| [`try-it.sh`](try-it.sh) | End-to-end tryout: creates a test collection, imports into it (one link invalid on purpose), reports the failure, then offers to delete it. |
| [`import-and-report.sh`](import-and-report.sh) `<links.json> <collection-id>` | Imports links and lists the ones that failed. Exits 1 if any failed, so it fits in cron or CI. |
| [`export-links-csv.sh`](export-links-csv.sh) `[output.csv] [collection-id]` | Exports links to a CSV spreadsheet. |
| [`collection-urls.sh`](collection-urls.sh) `<collection-id>` | Prints every URL in a collection, one per line. |
| [`link-details.sh`](link-details.sh) `<link-id>` | Shows a link's description and dates, which the table leaves out. |
| [`rearchive-collection.sh`](rearchive-collection.sh) `<collection-id>` | Pipes a collection's link IDs from `-q` into `links archive`, to refresh every snapshot. |

`try-it.sh` changes your account: it creates a collection called **CLI test** plus the date and time and,
if you answer `y`, deletes it and its links at the end. If it stops early, it prints the command to delete it.

All scripts check their arguments and that `linkwarden` and `jq` are installed before calling the API.
They exit 0 on success, 1 when something failed and 2 on a usage mistake. Shared helpers live in
[`lib.sh`](lib.sh).

The scripts follow the [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html):
each starts with a header that says what it does and how to call it, and every function has a comment
block.

## One-liners

```bash
# See what --json gives you
linkwarden collections list --json | jq

# Bare IDs, one per line
linkwarden collections list -q

# Capture the ID of something you create
COL=$(linkwarden collections create --name "Reading List" -q)
```

## Good to know

- Lists print as a JSON array (empty is `[]`), a single `get`, `create` or `update` as an object. Field
  names are the API's own.
- `links list` returns only the first page of links, so exports of large libraries are incomplete.
- Errors and info messages go to stderr, so stdout is always clean JSON or IDs. Check the exit code: it's
  non-zero on any failure, including a `links import` where some links failed.
