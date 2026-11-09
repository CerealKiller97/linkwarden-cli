#!/usr/bin/env bash
# Print every URL in a collection, one per line, e.g. to open them all or check
# them with curl.
#
# Usage:
#   examples/collection-urls.sh <collection-id>
#
# Pick an ID from: linkwarden collections list
set -euo pipefail
. "$(dirname "$0")/lib.sh"

[ $# -eq 1 ] || usage "<collection-id>"
require_id "collection-id" "$1"

"$LW" links list --collection-id "$1" --json | jq -r '.[].url'
