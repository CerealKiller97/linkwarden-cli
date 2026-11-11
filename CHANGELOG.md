# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-11-11

The first release of `linkwarden`, a command-line client for [Linkwarden](https://linkwarden.app). It works
with Linkwarden Cloud and self-hosted instances, on macOS and Linux (amd64 and arm64).

### ✨ Added

- **links:** List, get, create, update and delete links, with bulk update and bulk delete. `links create <url>`
  saves one link with an optional name, description, collection and tags.
- **links:** `links import --file` saves many links from a JSON file behind a progress bar. Links already in the
  collection are skipped, `--concurrency` sets how many are saved at once, and the command exits non-zero if any
  link fails.
- **links:** Archive a link on demand, delete its archives, and list its highlights.
- **collections:** List, get, create, update and delete collections, including nested collections.
- **tags:** List, get, create, update and delete tags, with bulk delete.
- **search:** Search your links, with `--sort` by date, name or description.
- **archives:** Download a link's PNG, JPEG, PDF, HTML or readable snapshot, and upload your own.
- **highlights, rss, dashboard:** Manage highlights and RSS subscriptions, and view and update the dashboard.
- **tokens:** List, generate and revoke API tokens.
- **migrations:** Export your library, and import from Linkwarden, HTML bookmarks, Wallabag, Omnivore or Pocket.
- **config, favicon:** Show the instance's configuration and fetch a site's favicon.
- **init:** `linkwarden init` signs you in with an access token, email and password, or any SSO provider the
  instance offers, checks it, and saves it to `$XDG_CONFIG_HOME/linkwarden/config.json`. `LINKWARDEN_TOKEN` and
  `LINKWARDEN_URL` override the file.
- **interactive mode:** In a terminal, missing arguments open filterable, paged pickers (the link picker loads
  more links as you scroll), results show as tables with clickable URLs, a spinner shows while the API works, and
  deletes and revokes ask No / Yes first.
- **scripted mode:** `--json` prints results as JSON, `-q` / `--quiet` prints only IDs, and `--non-interactive`
  (or no terminal) fails fast instead of prompting. Deletes need `--yes` there. Logs go to stderr, and Ctrl+C
  cancels in-flight requests and exits with status 130.
- **install:** Install with a one-line `curl | bash` script (with shell completions), Homebrew
  (`brew install CerealKiller97/tap/linkwarden-cli`), `go install`, a release archive, or the multi-arch Docker
  image `ghcr.io/cerealkiller97/linkwarden-cli`.
- **completions:** Shell completions for bash, zsh and fish.

[Unreleased]: https://github.com/CerealKiller97/linkwarden-cli/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/CerealKiller97/linkwarden-cli/releases/tag/v1.0.0
