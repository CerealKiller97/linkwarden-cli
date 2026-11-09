#!/bin/bash
#
# End-to-end tryout of the --json and -q flags against your account.
#
# Companion to the other scripts here: creates a test collection and keeps its
# ID with -q, imports two links into it (one invalid on purpose) through
# import-and-report.sh, lists what was saved with --json, then asks before
# deleting the collection and its links. The collection is named "CLI test"
# plus the date and time, so reruns don't clash. If the script stops early (an
# error, or Ctrl+C), it prints the command to delete what it created.
#
# Usage: examples/try-it.sh
#   LINKWARDEN_BIN=bin/linkwarden examples/try-it.sh   # use a local build

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
readonly SCRIPT_DIR
source "${SCRIPT_DIR}/lib.sh"

# State the EXIT trap needs, so it lives outside main.
collection_id=""
cleaned_up=false
links_file=""

#######################################
# Prints a step heading in bold.
# Arguments:
#   The heading.
# Outputs:
#   Writes the heading to STDOUT, after a blank line.
#######################################
step() {
  printf '\n\033[1m%s\033[0m\n' "$*"
}

#######################################
# Removes the temp file and, if the script stopped before step 5, says how to
# delete the test collection. Runs on EXIT.
# Globals:
#   LW
#   collection_id
#   cleaned_up
#   links_file
# Outputs:
#   Writes the delete command to STDOUT when the collection was left behind.
#######################################
cleanup() {
  if [[ -n "${links_file}" ]]; then
    rm -f "${links_file}"
  fi
  if [[ -n "${collection_id}" && "${cleaned_up}" == false ]]; then
    echo
    echo "The test collection ${collection_id} is still there. Delete it with:"
    echo "  ${LW} collections delete --id ${collection_id}"
  fi
}

#######################################
# Runs the tryout, one step at a time.
# Globals:
#   LW
#   SCRIPT_DIR
#   collection_id
#   cleaned_up
#   links_file
# Outputs:
#   Writes each step and its result to STDOUT.
#######################################
main() {
  check_deps
  trap cleanup EXIT

  local name
  name="CLI test $(date '+%Y-%m-%d %H.%M.%S')"

  step "1. Create a collection, keeping its ID with -q"
  collection_id="$("${LW}" collections create --name "${name}" -q)"
  echo "Created \"${name}\" with ID ${collection_id}"

  step "2-3. Import two links, one invalid, and report the failure"
  links_file="$(mktemp)"
  cat > "${links_file}" <<'EOF'
[
  { "url": "https://go.dev/doc/effective_go", "title": "Effective Go" },
  { "url": "not a url", "title": "Invalid on purpose" }
]
EOF
  local status=0
  "${SCRIPT_DIR}/import-and-report.sh" "${links_file}" "${collection_id}" \
    || status=$?
  if (( status == 0 )); then
    echo "(No failures this time.)"
  else
    echo "import-and-report.sh exited ${status}, as a cron job would see it."
  fi

  step "4. Links now in the collection (--json, shaped with jq)"
  "${LW}" links list --collection-id "${collection_id}" --json \
    | jq -r '.[] | "  \(.id)  \(.name)  \(.url)"'

  step "5. Clean up"
  local prompt="Delete collection ${collection_id} and its links?"
  prompt+=" This can't be undone. [y/N] "
  local answer
  read -r -p "${prompt}" answer || answer=""
  case "${answer}" in
    y | Y)
      "${LW}" collections delete --id "${collection_id}" --yes -q
      echo "Deleted."
      ;;
    *)
      echo "Kept it. Delete it later with:"
      echo "  ${LW} collections delete --id ${collection_id}"
      ;;
  esac
  cleaned_up=true
}

main "$@"
