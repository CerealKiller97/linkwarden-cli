#!/usr/bin/env bash
# Import links from a JSON file and report any that failed. Exits 1 when a link
# failed, so it fits in cron or CI.
#
# Usage:
#   examples/import-and-report.sh <links.json> <collection-id>
#
# The file is a JSON array of {"url": "...", "title": "..."}.
set -uo pipefail
. "$(dirname "$0")/lib.sh"

[ $# -eq 2 ] || usage "<links.json> <collection-id>"
file="$1"
collection="$2"
[ -r "$file" ] || die "can't read $file"
require_id "collection-id" "$collection"

# The CLI logs each failed link on stderr as well. Keep that aside: the report
# below lists the same failures more readably, and the log is only shown when
# the import stops before it gets to a summary.
errlog=$(mktemp)
trap 'rm -f "$errlog"' EXIT

report=$("$LW" --non-interactive links create --file "$file" --collection "$collection" --json 2>"$errlog")
status=$?

if [ -z "$report" ]; then
  # Bad file, unknown collection, server unreachable… The CLI said why.
  cat "$errlog" >&2
  exit "$status"
fi

echo "$report" | jq -r '"\(.total) links: \(.success) created, \(.alreadyExists) already saved, \(.failed) failed"'

if [ "$(echo "$report" | jq '.failed')" -gt 0 ]; then
  echo "Failed:" >&2
  echo "$report" | jq -r '.failures[] | "  \(.url)  \(.error | sub("^linkwarden: "; ""))"' >&2
fi

exit "$status"
