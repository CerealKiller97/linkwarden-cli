#!/usr/bin/env bash
# Export links to a CSV file you can open in a spreadsheet: id, name, url, collection.
#
# Usage:
#   examples/export-links-csv.sh [output.csv] [collection-id]
#
# Without a collection ID it exports from all collections. Only the first page
# of links is exported (links list doesn't page through everything yet).
set -euo pipefail
. "$(dirname "$0")/lib.sh"

output="${1:-links.csv}"
collection="${2:-}"

args=(links list --json)
if [ -n "$collection" ]; then
  require_id "collection-id" "$collection"
  args+=(--collection-id "$collection")
fi

# Write to a temp file and move it into place, so a failed fetch doesn't leave
# a half-written CSV behind.
tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
{
  echo '"id","name","url","collection"'
  "$LW" "${args[@]}" | jq -r '.[] | [.id, .name, .url, .collection.name] | @csv'
} > "$tmp"
mv "$tmp" "$output"

echo "Wrote $(($(wc -l < "$output") - 1)) links to $output" >&2
