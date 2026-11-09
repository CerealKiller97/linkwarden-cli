# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### ♻️ Changed

- **ui:** Replace promptui with Bubble Tea for all interactive prompts. Pickers can now be filtered by typing, scroll past 10 items, and cancel with Esc.
- **ui:** Replace go-pretty tables with Lip Gloss tables: rounded border, bold header, title badge and dimmed alternate rows. The `links import` summary colors its counts (blue / amber / red). In a terminal, tables taller than the screen open as an inline, scrollable bubbles/table that stays on screen when you quit.
- **ui:** Replace the schollz/progressbar bar in `links import` with a bubbles/progress bar on stderr (hidden when stderr isn't a terminal). Per-link failures are logged after the bar finishes instead of interleaving with it.
- **tokens:** `tokens generate` and `tokens revoke` now fail fast without an interactive terminal, like the other pickers.
- **links:** The JSON-file import moved from `links create --file` to `links import --file`. `links create --file` still works but prints a deprecation notice.
- **cli:** ⚠️ **Breaking:** commands that delete or revoke now refuse to run without an interactive terminal (or with `--non-interactive`) unless you pass `--yes`. Add `--yes` to scripts and cron jobs that delete things.

### ✨ Added

- **links:** `links create <url>` saves a single link, with optional `--name`, `--description`, `--collection` (defaults to Unorganized) and repeatable `--tag`.
- **cli:** Global `--json` flag prints results as JSON with the API's field names instead of a table: lists as arrays, a single `get` or `create` as an object, and the `links import` summary with each failed URL and why. Global `-q` / `--quiet` prints only the ID of each result, one per line, and hides info messages, e.g. `linkwarden tags list -q | xargs …`.
- **cli:** `links delete`, `links bulk-delete`, `links delete-archive`, `collections delete`, `tags delete`, `tags bulk-delete`, `rss delete`, `highlights delete` and `tokens revoke` ask **No / Yes** before acting, with No as the default answer. `collections delete` names the collection and how many links go with it. `--yes` / `-y` skips the question.
- **init:** Add `linkwarden init`, a guided setup that signs in with an access token (opening the Access Tokens page in your browser), email/password, or any SSO provider the instance offers (Google, Apple, …), verifies it and saves it to `$XDG_CONFIG_HOME/linkwarden/config.json`. Commands fall back to this file when `LINKWARDEN_TOKEN` isn't set.
- **ui:** All prompts, tables and the progress bar use Linkwarden's own blue palette (`#8FD1F8` / `#4AA3E2`, defined once in `pkg/tui/theme.go`).
- **ui:** URLs in the links table (`links list`, `links get`, `search`, `dashboard`, `archives`) and the RSS table are clickable in terminals that support hyperlinks (iTerm2, Ghostty, kitty, WezTerm, VS Code, Windows Terminal, …). A URL too long for the column is still shown cut short with `…`, but clicking it opens the whole URL. Piped output is unchanged.
- **ui:** Pickers page with the bubbles paginator (←/→ or PgUp/PgDn change page). The link picker (`links get`, `links delete`, `archives get`, `highlights create`, …) shows each link's shortened URL next to its name, and loads the next page of links from the API as you reach the last one loaded, or while a filter has fewer matches than fit on a page. Before, it only ever showed the newest 50 links.
- **links:** `links import` skips links already in the target collection (and URLs repeated in the file), and reports them in a new **Already exists** column. A 409 from accounts with "prevent duplicate links" on is counted there too.

### 🐛 Fixed

- **cli:** `links create` and `tokens generate` now have a description in `--help`.
- **cli:** Trim trailing slashes from `LINKWARDEN_URL`.
- **tokens:** `tokens revoke` revokes the selected token even when several tokens share a name.

## [1.0.0] - 2026-10-06

### ✨ Added

- **links:** Collection option added, if nothing passed show interactive prompt for collections.
- **sdk:** Add Linkwarden API client and full CRUD commands for links, collections, tags, tokens.
- **sdk:** Add Search/RSS/Highlights/Dashboard SDK services and CLI commands.
- **sdk:** Add Archives/Config/Favicon SDK services, implement migrations import.
- **ui:** Replace forced-color table style with adaptive rounded + cyan header.
- **cli:** Add global --non-interactive flag to force fail-fast prompting.
- **tokens:** Show colored manual setup steps when user skips auto env config.

### ♻️ Changed

- Binary/command name from linkwarden-cli to linkwarden.

### 🐛 Fixed

- **utils:** Correct shell config paths for token persistence, drop broken reload.
- **utils:** Fail fast instead of prompting when stdin isn't interactive.

[Unreleased]: https://github.com/CerealKiller97/linkwarden-cli/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/CerealKiller97/linkwarden-cli/releases/tag/v1.0.0
