#!/bin/bash
#
# Exports links to a CSV file you can open in a spreadsheet.
#
# Companion to `linkwarden links list --json`: writes one row per link with its
# id, name, url and collection, under a header row. Without a collection ID it
# exports from all collections. The file is written in one go at the end, so a
# failed fetch never leaves a half-written CSV behind. Only the first page of
# links is exported, since `links list` doesn't page through everything yet.
#
# Usage: examples/export-links-csv.sh [output.csv] [collection-id]
#   examples/export-links-csv.sh                    # all links to links.csv
#   examples/export-links-csv.sh reading.csv 7      # one collection
#   LINKWARDEN_BIN=bin/linkwarden examples/export-links-csv.sh

set -euo pipefail

source "$(dirname "$0")/lib.sh"

#######################################
# Fetches the links and writes the CSV.
# Globals:
#   LW
# Arguments:
#   Path of the CSV to write; defaults to links.csv.
#   ID of a collection to export; optional.
# Outputs:
#   Writes the CSV file, and how many links it holds to STDERR.
#######################################
main() {
  local output="${1:-links.csv}"
  local collection="${2:-}"
  local -a args=(links list --json)

  if [[ -n "${collection}" ]]; then
    require_id "collection-id" "${collection}"
    args+=(--collection-id "${collection}")
  fi
  check_deps

  local tmp
  tmp="$(mktemp)"
  trap "rm -f '${tmp}'" EXIT

  {
    echo '"id","name","url","collection"'
    "${LW}" "${args[@]}" \
      | jq -r '.[] | [.id, .name, .url, .collection.name] | @csv'
  } > "${tmp}"
  mv "${tmp}" "${output}"

  local lines
  lines="$(wc -l < "${output}")"
  echo "Wrote $(( lines - 1 )) links to ${output}" >&2
}

main "$@"
