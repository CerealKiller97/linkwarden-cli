# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### ♻️ Changed

- **ui:** Replace promptui with Bubble Tea for all interactive prompts. Pickers can now be filtered by typing, scroll past 10 items, and cancel with Esc.
- **tokens:** `tokens generate` and `tokens revoke` now fail fast without an interactive terminal, like the other pickers.

### 🐛 Fixed

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
