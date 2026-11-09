#!/bin/bash
#
# Triggers a fresh archive of every link in a collection.
#
# Companion to `linkwarden links list -q` and `linkwarden links archive`: pipes
# the collection's link IDs into archive, one at a time, so Linkwarden takes a
# new screenshot, PDF and readable copy of each. Handy after changing archive
# settings, or when old snapshots are broken. Each link is best-effort, so one
# failure doesn't stop the rest; exits 1 if any failed. Only the first page of
# links is covered, since `links list` doesn't page through everything yet.
#
# Usage: examples/rearchive-collection.sh <collection-id>
#   LINKWARDEN_BIN=bin/linkwarden examples/rearchive-collection.sh 3

set -euo pipefail

source "$(dirname "$0")/lib.sh"

#######################################
# Re-archives each link in the collection.
# Globals:
#   LW
# Arguments:
#   ID of the collection.
# Outputs:
#   Writes how many links were archived to STDERR, and any errors.
# Returns:
#   0 when every link was archived, 1 otherwise.
#######################################
main() {
  (( $# == 1 )) || usage "<collection-id>"
  require_id "collection-id" "$1"
  check_deps

  # Fetch the IDs first: a failure inside `for id in $(...)` wouldn't stop
  # the script.
  local ids
  ids="$("${LW}" links list --collection-id "$1" -q)"

  local count=0
  local failed=0
  local id
  for id in ${ids}; do
    if "${LW}" links archive --id "${id}" -q; then
      count=$(( count + 1 ))
    else
      failed=$(( failed + 1 ))
    fi
  done

  if (( failed > 0 )); then
    echo "Archiving triggered for ${count} links; ${failed} failed" \
      "(see above)." >&2
    return 1
  fi
  echo "Archiving triggered for ${count} links." >&2
}

main "$@"
