#!/usr/bin/env bash
#
# Smoke-tests a built CLI binary against a live Linkwarden instance.
#
#   scripts/e2e.sh [path/to/linkwarden]   (defaults to bin/linkwarden)
#
# Reads LINKWARDEN_URL and LINKWARDEN_TOKEN from the environment. Locally:
#   scripts/dev-linkwarden.sh up
#   eval "$(scripts/dev-linkwarden.sh env)"
#   go build -o bin/linkwarden . && scripts/e2e.sh
#
# It creates a collection, tags and links, so point it at a throwaway instance only.

set -euo pipefail

BIN="${1:-bin/linkwarden}"

: "${LINKWARDEN_URL:?LINKWARDEN_URL is not set}"
: "${LINKWARDEN_TOKEN:?LINKWARDEN_TOKEN is not set}"

WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT

# lw ARGS... — runs the CLI non-interactively, echoes the command, and keeps its output in $OUT.
lw() {
	echo "▶ linkwarden $*"
	OUT="$("$BIN" --non-interactive "$@" </dev/null 2>&1)" || {
		echo "$OUT"
		echo "✖ linkwarden $* failed" >&2
		exit 1
	}
	echo "$OUT"
}

# expect TEXT — fails unless the last command's output contains TEXT.
expect() {
	grep -qF -- "$1" <<<"$OUT" || {
		echo "✖ expected output to contain: $1" >&2
		exit 1
	}
}

lw --version
lw config get

lw collections create --name "E2E Smoke"
collection_id="$(grep -oE 'id=[0-9]+' <<<"$OUT" | head -n 1 | cut -d= -f2)"
[[ -n "$collection_id" ]] || { echo "✖ could not read the new collection's ID" >&2; exit 1; }

lw collections list
expect "E2E Smoke"

lw tags create --label e2e --label smoke
expect "e2e"

cat >"$WORK_DIR/links.json" <<'EOF'
[
  {"url": "https://go.dev/doc/effective_go", "title": "Effective Go"},
  {"url": "https://linkwarden.app", "title": "Linkwarden"}
]
EOF
lw links create --file "$WORK_DIR/links.json" --collection "$collection_id"

lw links list --collection-id "$collection_id"
expect "https://go.dev/doc/effective_go"
expect "https://linkwarden.app"

# Links reach the search index asynchronously, so retry for up to a minute.
for attempt in $(seq 1 30); do
	if OUT="$("$BIN" --non-interactive search --query "Effective" </dev/null 2>&1)" &&
		grep -qF "https://go.dev/doc/effective_go" <<<"$OUT"; then
		echo "▶ linkwarden search --query Effective"
		echo "$OUT"
		break
	fi
	if [[ "$attempt" == 30 ]]; then
		echo "$OUT"
		echo "✖ search never returned the new link" >&2
		exit 1
	fi
	sleep 2
done

echo "✔ e2e smoke test passed"
