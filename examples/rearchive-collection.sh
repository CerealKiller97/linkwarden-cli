#!/usr/bin/env bash
# Pipe IDs from one command into another: list a collection's link IDs with -q
# and trigger a fresh archive (screenshot, PDF, readable copy) of each. Handy
# after changing archive settings, or when old snapshots are broken.
#
# Usage:
#   examples/rearchive-collection.sh <collection-id>
set -euo pipefail
. "$(dirname "$0")/lib.sh"

[ $# -eq 1 ] || usage "<collection-id>"
require_id "collection-id" "$1"

# Fetch the IDs first: a failure inside `for id in $(...)` wouldn't stop the script.
ids=$("$LW" links list --collection-id "$1" -q)

count=0
failed=0
for id in $ids; do
  if "$LW" links archive --id "$id" -q; then
    count=$((count + 1))
  else
    failed=$((failed + 1))
  fi
done

if [ "$failed" -gt 0 ]; then
  echo "Archiving triggered for $count links; $failed failed (see above)." >&2
  exit 1
fi
echo "Archiving triggered for $count links." >&2
