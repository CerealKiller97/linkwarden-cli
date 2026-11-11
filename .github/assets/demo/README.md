# Demo recordings

The GIFs in the main README are recorded with [VHS](https://github.com/charmbracelet/vhs) from the tapes in this folder.

| Tape | Shows |
| --- | --- |
| [`interactive.tape`](interactive.tape) | Tables, the link picker and the delete confirmation |
| [`collections.tape`](collections.tape) | Creating and nesting collections, updating one through a picker, deleting with confirmation |
| [`links.tape`](links.tape) | Saving, listing, searching, bulk-tagging and renaming links |
| [`import.tape`](import.tape) | `links import` with a collection picker, progress bar and results table |
| [`scripting.tape`](scripting.tape) | `--json`, `-q` and failing fast without a terminal |
| [`../install.tape`](../install.tape) | The install script (it has its own setup; see the tape) |

Each tape writes `.github/assets/<name>.gif`.

## Recording

```bash
.github/assets/demo/record.sh                    # every tape
.github/assets/demo/record.sh collections links  # only these
```

[`record.sh`](record.sh) needs Docker, Go, `curl`, `jq` and `openssl`. It never touches your own account:

1. Starts a throwaway Linkwarden on port 3200 from [`compose.yml`](compose.yml), with fresh volumes.
2. Creates a demo user and API token, and fills the account with [`seed.sh`](seed.sh).
3. Builds the CLI for Linux into `bin/` (git-ignored), plus a VHS image that adds `curl` and `jq`.
4. Records each tape into `.github/assets/`.
5. Removes the Linkwarden containers and volumes.

| Variable | Default | Effect |
| --- | --- | --- |
| `DEMO_PORT` | `3200` | Port for the throwaway Linkwarden |
| `LINKWARDEN_VERSION` | `latest` | Linkwarden image tag |
| `KEEP` | unset | Set to leave Linkwarden running afterwards, to try commands against the demo data |

## Writing a tape

- Start with `Source .github/assets/demo/setup.tape`. It sets the theme, size and font, and hides the shell setup in
  [`bashrc`](bashrc) (prompt, `linkwarden` on `PATH`, working directory). Override `Set Width`, `Set Height` or
  `Set FontSize` after it when a table needs more room.
- Type a short `# comment` before each step. Bash ignores it, and it tells the viewer what's coming.
- The tapes rely on the IDs that `seed.sh` produces on an empty account (collection `Go` is 1, link
  `Go Concurrency Patterns: Pipelines` is 2) and run in the order listed in `record.sh`.
- Clean up between `Hide` and the end of the tape: delete what you created and undo what you changed, so the next
  tape sees the seed data. Inside `xargs`, call `linkwarden-demo`: `linkwarden` is a shell function there.
- To check a frame without watching the GIF, add `Screenshot bin/check.png` at that point while you work on it.
