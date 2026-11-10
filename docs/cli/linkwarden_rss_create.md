## linkwarden rss create

Create an RSS subscription

```bash
linkwarden rss create [flags]
```

### Examples

```bash
# Follow a feed into collection 3
linkwarden rss create --name "The Go Blog" --url https://go.dev/blog/feed.atom --collection-id 3

# Name the destination collection instead
linkwarden rss create --name "The Go Blog" --url https://go.dev/blog/feed.atom --collection-name Go
```

### Options

```
      --collection-id int        Destination collection ID.
      --collection-name string   Destination collection name (used if collection-id is not set).
  -h, --help[=true|false]        help for create
      --name string              Subscription name.
      --url string               Feed URL.
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden rss](linkwarden_rss.md)	 - Manage RSS subscriptions
