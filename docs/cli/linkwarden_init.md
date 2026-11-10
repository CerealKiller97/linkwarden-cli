## linkwarden init

Set up the CLI with your Linkwarden account

### Synopsis

Walks you through signing in to your Linkwarden instance, with an access token (recommended), your email and password, or a provider such as Google or Apple, and saves the result to the CLI's config file.

```bash
linkwarden init [flags]
```

### Examples

```bash
# Sign in (run it again any time to switch accounts or instances)
linkwarden init
```

### Options

```
  -h, --help[=true|false]   help for init
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden](linkwarden.md)	 - Linkwarden is a CLI for the Linkwarden platform
