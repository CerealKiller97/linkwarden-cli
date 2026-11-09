#!/usr/bin/env bash
# Show fields of one link that the table doesn't: description and dates.
#
# Usage:
#   examples/link-details.sh <link-id>
set -euo pipefail
. "$(dirname "$0")/lib.sh"

[ $# -eq 1 ] || usage "<link-id>"
require_id "link-id" "$1"

"$LW" links get --id "$1" --json | jq '{name, url, description, createdAt}'
