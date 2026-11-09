#!/bin/bash
#
# Imports links from a JSON file into a collection and reports any that failed.
#
# Companion to `linkwarden links import --json`: prints a one-line summary of
# the import, then each failed URL with the server's reason. The CLI's own
# per-link log lines are held back, since the report lists the same failures;
# they're shown only when the import stops before it reaches a summary (bad
# file, unknown collection, server down). Exits 1 when any link failed, so it
# fits in cron or CI. The file is a JSON array of objects with "url" and
# "title".
#
# Usage: examples/import-and-report.sh <links.json> <collection-id>
#   examples/import-and-report.sh links.json 3 || notify-send "Import failed"
#   LINKWARDEN_BIN=bin/linkwarden examples/import-and-report.sh links.json 3

set -euo pipefail

source "$(dirname "$0")/lib.sh"

#######################################
# Imports the file and prints the report.
# Globals:
#   LW
# Arguments:
#   Path to the JSON file of links.
#   ID of the collection to import into.
# Outputs:
#   Writes the summary to STDOUT, and the failed links to STDERR.
# Returns:
#   0 when every link was created or already saved, non-zero otherwise.
#######################################
main() {
  (( $# == 2 )) || usage "<links.json> <collection-id>"
  local file="$1"
  local collection="$2"
  [[ -r "${file}" ]] || die "can't read ${file}"
  require_id "collection-id" "${collection}"
  check_deps

  local errlog
  errlog="$(mktemp)"
  # Expand the path now, while the local is in scope; the trap runs after
  # main returns.
  # shellcheck disable=SC2064
  trap "rm -f '${errlog}'" EXIT

  local report
  local status=0
  report="$("${LW}" --non-interactive links import --json \
    --file "${file}" --collection "${collection}" 2> "${errlog}")" \
    || status=$?

  if [[ -z "${report}" ]]; then
    cat "${errlog}" >&2
    return "${status}"
  fi

  echo "${report}" | jq -r '
    "\(.total) links: \(.success) created, "
    + "\(.alreadyExists) already saved, \(.failed) failed"'

  if (( $(echo "${report}" | jq '.failed') > 0 )); then
    echo "Failed:" >&2
    echo "${report}" \
      | jq -r '.failures[]
          | "  \(.url)  \(.error | sub("^linkwarden: "; ""))"' >&2
  fi

  return "${status}"
}

main "$@"
