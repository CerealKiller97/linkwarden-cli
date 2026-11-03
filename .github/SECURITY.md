# Security Policy

## Supported versions

Only the latest release of linkwarden-cli receives security fixes. Please upgrade to the
[newest release](https://github.com/CerealKiller97/linkwarden-cli/releases/latest) before reporting.

| Version          | Supported |
| ---------------- | --------- |
| Latest release   | ✅        |
| Older releases   | ❌        |

## Reporting a vulnerability

**Please do not open a public issue, discussion or pull request for security problems.**

Report vulnerabilities privately through GitHub's
[private vulnerability reporting](https://github.com/CerealKiller97/linkwarden-cli/security/advisories/new).

Please include:

- The affected version (`linkwarden --version`) or commit, and your OS
- A description of the issue and its impact
- Steps to reproduce, or a proof of concept
- Any suggested fix, if you have one

## What to expect

- An acknowledgement within **7 days**.
- An assessment and, if confirmed, a plan for a fix within **30 days**.
- A patched release and a published GitHub Security Advisory, crediting you unless you
  prefer to stay anonymous.

## Scope

In scope:

- The `linkwarden` CLI. Report issues in the Go SDK to
  [linkwarden-go-sdk](https://github.com/CerealKiller97/linkwarden-go-sdk/security)
- Handling of API tokens, for example leaking `LINKWARDEN_TOKEN` in logs or errors, or
  writing it to shell config files with unsafe permissions
- Release artifacts and the release workflow

Out of scope:

- Vulnerabilities in Linkwarden itself (the server or web app). Report those to the
  [Linkwarden project](https://github.com/linkwarden/linkwarden/security).
- The throwaway local instance started by `scripts/dev-linkwarden.sh`, which is meant
  for development only and uses weak default credentials.
