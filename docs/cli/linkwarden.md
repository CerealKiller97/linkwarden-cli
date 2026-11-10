## linkwarden

Linkwarden is a CLI for the Linkwarden platform

### Synopsis

Manage your Linkwarden bookmarks from the terminal: links, collections, tags, archives,
RSS feeds, highlights and backups, on Linkwarden Cloud or a self-hosted instance.

Run "linkwarden init" to sign in. Leave out --id on most commands to pick from a list,
add --json for machine-readable output, or -q for bare IDs.

### Options

```
  -h, --help[=true|false]              help for linkwarden
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden archives](linkwarden_archives.md)	 - Manage link archive files
* [linkwarden collections](linkwarden_collections.md)	 - Manage collections
* [linkwarden completion](linkwarden_completion.md)	 - Generate the autocompletion script for the specified shell
* [linkwarden config](linkwarden_config.md)	 - View instance configuration
* [linkwarden dashboard](linkwarden_dashboard.md)	 - View and configure the dashboard
* [linkwarden favicon](linkwarden_favicon.md)	 - Fetch favicons
* [linkwarden highlights](linkwarden_highlights.md)	 - Manage link highlights
* [linkwarden init](linkwarden_init.md)	 - Set up the CLI with your Linkwarden account
* [linkwarden links](linkwarden_links.md)	 - Manage links
* [linkwarden migrations](linkwarden_migrations.md)	 - Manage migration (Import / Export)
* [linkwarden rss](linkwarden_rss.md)	 - Manage RSS subscriptions
* [linkwarden search](linkwarden_search.md)	 - Search links
* [linkwarden tags](linkwarden_tags.md)	 - Manage tags
* [linkwarden tokens](linkwarden_tokens.md)	 - Manage tokens
