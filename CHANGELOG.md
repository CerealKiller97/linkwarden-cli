# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### ♻️ Changed

- **sdk:** `WithLogger` now takes a standard-library `*slog.Logger` instead of `zerolog.Logger`, so the SDK no longer depends on zerolog. **Breaking** for SDK users who passed a logger.
- **sdk:** Without `WithLogger`, the client now logs to stdout via a text `slog` handler instead of discarding logs.
- **cli:** Replace zerolog with `log/slog` and [tint](https://github.com/lmittmann/tint) for colorized console logs. Same fields (`id`, `ids`, `path`, `error`, …) and levels; color is now disabled when output isn't a terminal.

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
