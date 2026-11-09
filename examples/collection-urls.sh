#!/bin/bash
#
# Prints every URL in a collection, one per line.
#
# Companion to `linkwarden links list --json`: handy for feeding a collection
# into other tools, such as opening every link or checking them with curl.
# Pick a collection ID from `linkwarden collections list`. Only the first page
# of links is printed, since `links list` doesn't page through everything yet.
#
# Usage: examples/collection-urls.sh <collection-id>
#   examples/collection-urls.sh 3 | xargs -n1 curl -sIo /dev/null \
#     -w '%{http_code} %{url_effective}\n'          # spot dead links
#   LINKWARDEN_BIN=bin/linkwarden examples/collection-urls.sh 3

set -euo pipefail

source "$(dirname "$0")/lib.sh"

#######################################
# Prints the collection's URLs.
# Globals:
#   LW
# Arguments:
#   ID of the collection.
# Outputs:
#   Writes one URL per line to STDOUT.
#######################################
main() {
  (( $# == 1 )) || usage "<collection-id>"
  require_id "collection-id" "$1"
  check_deps

  "${LW}" links list --collection-id "$1" --json | jq -r '.[].url'
}

main "$@"
