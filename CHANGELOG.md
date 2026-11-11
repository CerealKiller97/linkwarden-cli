# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-11-11

### ✨ Added

- **links:** Collection option added, if nothing passed show interactive prompt for collections.
- **sdk:** Add Linkwarden API client and full CRUD commands for links, collections, tags, tokens.
- **sdk:** Add Search/RSS/Highlights/Dashboard SDK services and CLI commands.
- **sdk:** Add Archives/Config/Favicon SDK services, implement migrations import.
- **ui:** Replace forced-color table style with adaptive rounded + cyan header.
- **cli:** Add global --non-interactive flag to force fail-fast prompting.
- **tokens:** Show colored manual setup steps when user skips auto env config.
- SECURITY.md.
- CODEOWNERS.
- Update CODEOWNERS file.
- **release:** Publish to Homebrew, AUR and winget with GoReleaser.
- **release:** Publish multi-arch Docker image to GHCR.
- **links:** Create links concurrently with --concurrency flag.
- **links:** Skip existing links in create; restyle tables.
- **ui:** Use Linkwarden's blue palette for every TUI component.
- **init:** Add guided `linkwarden init` setup.
- **ui:** Standardize every command's result table.
- **ui:** Show a spinner while waiting on the Linkwarden API.
- **cli:** Confirm before deleting or revoking, with --yes to skip.
- **ui:** Make table URLs clickable, opening the full URL.
- **ui:** Page the link picker and show each link's short URL.
- **ui:** Keep the spinner up for at least a second once it appears.
- **cli:** Add --json and --quiet output flags.
- Add a one-line curl | bash installer.
- **install:** Default to /usr/local/bin and always install completions.
- **links:** Add links import and a single-link links create.
- **links:** Drop --file from links create.
- **cli:** Take names for --format and --sort.
- **release:** Add a Docker section to the release notes.
- **release:** Add missing winget manifest metadata.

### ♻️ Changed

- Binary/command name from linkwarden-cli to linkwarden.
- Replace pkg/linkwarden with linkwarden-go-sdk.
- Replace promptui with Bubble Tea for interactive prompts.
- Move tables and progress bar to Bubbles/Lip Gloss.
- **cli:** Use consistent int flags for IDs.
- Order struct fields for tighter memory alignment.

### 🐛 Fixed

- **utils:** Correct shell config paths for token persistence, drop broken reload.
- **utils:** Fail fast instead of prompting when stdin isn't interactive.
- **xdg:** Resolve XDG base directories consistently on Linux, macOS and Windows.
- **links:** Reject unknown --collection ID in links create.
- **ui:** Style tables after Bubble Tea's table example.
- **cli:** Trim trailing slashes from LINKWARDEN_URL.
- **utils:** Validate URLs before opening the browser (gosec G204).
- **cli:** Keep logs off stdout and report each error once.
- **cli:** Log the final command error through zerolog.
- **tui:** Pre-fill default inputs instead of using a placeholder.
- **collections:** Cache the Unorganized collection ID per instance.
- Exit quietly when a prompt is cancelled.
- **links:** Update the Linkwarden SDK to v1.0.1 for the search route.
- **links:** Exit non-zero when links create has failures.
- **cli:** Cancel in-flight requests on Ctrl+C.
- **cli:** Name the go install binary linkwarden.

### ♻️ Changed

- **ui:** Replace promptui with Bubble Tea for all interactive prompts. Pickers can now be filtered by typing, scroll past 10 items, and cancel with Esc.
- **ui:** Replace go-pretty tables with Lip Gloss tables: rounded border, bold header, title badge and dimmed alternate rows. The `links import` summary colors its counts (blue / amber / red). In a terminal, tables taller than the screen open as an inline, scrollable bubbles/table that stays on screen when you quit.
- **ui:** Replace the schollz/progressbar bar in `links import` with a bubbles/progress bar on stderr (hidden when stderr isn't a terminal). Per-link failures are logged after the bar finishes instead of interleaving with it.
- **tokens:** `tokens generate` and `tokens revoke` now fail fast without an interactive terminal, like the other pickers.
- **links:** the JSON-file import moved from `links create --file` to `links import --file`. `links create` no longer takes `--file`.
- **cli:** `--format` and `--sort` take names instead of the API's numbers: `archives get|upload|upload-standalone --format png|jpeg|pdf|html|readable`, `migrations import --format linkwarden|html|wallabag|omnivore|pocket` and `search --sort newest|oldest|name-asc|name-desc|description-asc|description-desc`. Numbers are rejected, so update scripts (e.g. `--format 2` → `--format pdf`). `--help` lists the names and shell completion suggests them. `archives upload` no longer offers `readable`, which is retrieve-only.
- **cli:** commands that delete or revoke now refuse to run without an interactive terminal (or with `--non-interactive`) unless you pass `--yes`. Add `--yes` to scripts and cron jobs that delete things.
- **cli:** ID flags are integers everywhere. `--id` and `--link-id` were strings and are now checked when the command line is parsed, like `--collection-id`, `--tag-id` and `--parent-id` already were, so a bad value gets the same error on every command; zero and negative IDs are rejected. `--ids` also takes repeated flags (`--ids 1,2 --ids 3`), and no longer accepts spaces inside the list (`--ids "1, 2"` → `--ids 1,2`) or an empty list.

### ✨ Added

- **links:** `links create <url>` saves a single link, with optional `--name`, `--description`, `--collection-id` (defaults to Unorganized) and repeatable `--tag`.
- **cli:** Global `--json` flag prints results as JSON with the API's field names instead of a table: lists as arrays, a single `get` or `create` as an object, and the `links import` summary with each failed URL and why. Global `-q` / `--quiet` prints only the ID of each result, one per line, and hides info messages, e.g. `linkwarden tags list -q | xargs …`.
- **cli:** `links delete`, `links bulk-delete`, `links delete-archive`, `collections delete`, `tags delete`, `tags bulk-delete`, `rss delete`, `highlights delete` and `tokens revoke` ask **No / Yes** before acting, with No as the default answer. `collections delete` names the collection and how many links go with it. `--yes` / `-y` skips the question.
- **init:** Add `linkwarden init`, a guided setup that signs in with an access token (opening the Access Tokens page in your browser), email/password, or any SSO provider the instance offers (Google, Apple, …), verifies it and saves it to `$XDG_CONFIG_HOME/linkwarden/config.json`. Commands fall back to this file when `LINKWARDEN_TOKEN` isn't set.
- **ui:** All prompts, tables and the progress bar use Linkwarden's own blue palette (`#8FD1F8` / `#4AA3E2`, defined once in `pkg/tui/theme.go`).
- **ui:** URLs in the links table (`links list`, `links get`, `search`, `dashboard`, `archives`) and the RSS table are clickable in terminals that support hyperlinks (iTerm2, Ghostty, kitty, WezTerm, VS Code, Windows Terminal, …). A URL too long for the column is still shown cut short with `…`, but clicking it opens the whole URL. Piped output is unchanged.
- **ui:** Pickers page with the bubbles paginator (←/→ or PgUp/PgDn change page). The link picker (`links get`, `links delete`, `archives get`, `highlights create`, …) shows each link's shortened URL next to its name, and loads the next page of links from the API as you reach the last one loaded, or while a filter has fewer matches than fit on a page. Before, it only ever showed the newest 50 links.
- **links:** `links import` skips links already in the target collection (and URLs repeated in the file), and reports them in a new **Already exists** column. A 409 from accounts with "prevent duplicate links" on is counted there too.

### 🐛 Fixed

- **cli:** Ctrl+C (or SIGTERM) now cancels in-flight API requests and exits with status 130 and a short "Interrupted" warning, instead of being cut off mid-request. `links import` starts no new requests once interrupted, still prints its summary of what was saved, and doesn't count the cancelled links as failures. A second Ctrl+C kills the process outright.
- **cli:** `links create` and `tokens generate` now have a description in `--help`.
- **cli:** Trim trailing slashes from `LINKWARDEN_URL`.
- **tokens:** `tokens revoke` revokes the selected token even when several tokens share a name.


[Unreleased]: https://github.com/CerealKiller97/linkwarden-cli/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/CerealKiller97/linkwarden-cli/releases/tag/v1.0.0
