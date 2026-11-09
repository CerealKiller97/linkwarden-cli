#!/bin/bash
#
# Shows the fields of one link that the links table leaves out.
#
# Companion to `linkwarden links get --json`: prints the link's name, URL,
# description and creation date as a JSON object. Change the jq filter below
# to pick other fields; `links get --id <id> --json` shows them all.
#
# Usage: examples/link-details.sh <link-id>
#   examples/link-details.sh 412 | jq -r .description   # just the description
#   LINKWARDEN_BIN=bin/linkwarden examples/link-details.sh 412

set -euo pipefail

source "$(dirname "$0")/lib.sh"

#######################################
# Prints the link's details.
# Globals:
#   LW
# Arguments:
#   ID of the link.
# Outputs:
#   Writes a JSON object to STDOUT.
#######################################
main() {
  (( $# == 1 )) || usage "<link-id>"
  require_id "link-id" "$1"
  check_deps

  "${LW}" links get --id "$1" --json \
    | jq '{name, url, description, createdAt}'
}

main "$@"
