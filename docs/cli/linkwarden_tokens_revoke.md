## linkwarden tokens revoke

Revoke a token

```bash
linkwarden tokens revoke [flags]
```

### Examples

```bash
# Pick the token from a list, then confirm
linkwarden tokens revoke

# Pick the token, and skip the confirmation
linkwarden tokens revoke --yes
```

### Options

```
  -h, --help[=true|false]   help for revoke
  -y, --yes[=true|false]    Skip the confirmation prompt (required when there's no interactive terminal).
```

### Options inherited from parent commands

```
      --json[=true|false]              Print results as JSON instead of a table.
      --non-interactive[=true|false]   Never prompt (e.g. for a picker); fail fast instead. On by default when stdin isn't a terminal.
  -q, --quiet[=true|false]             Print only the ID of each result, one per line, and hide info messages.
```

### SEE ALSO

* [linkwarden tokens](linkwarden_tokens.md)	 - Manage tokens
