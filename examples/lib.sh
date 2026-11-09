# shellcheck shell=bash
#
# Shared helpers for the example scripts in this folder.
#
# Finds the linkwarden binary, checks that the commands the scripts need are
# installed, and validates arguments, so a mistake is reported before any API
# call. It only defines functions and the LW constant; sourcing it runs
# nothing. Not meant to be run on its own.
#
# Usage: source "$(dirname "$0")/lib.sh"
#   LINKWARDEN_BIN=bin/linkwarden examples/try-it.sh   # use a local build

# The linkwarden binary to run: LINKWARDEN_BIN, or linkwarden from the PATH.
readonly LW="${LINKWARDEN_BIN:-linkwarden}"

#######################################
# Prints an error message.
# Arguments:
#   The message.
# Outputs:
#   Writes "error: <message>" to STDERR.
#######################################
err() {
  echo "error: $*" >&2
}

#######################################
# Prints an error message and exits with status 1.
# Arguments:
#   The message.
# Outputs:
#   Writes "error: <message>" to STDERR.
#######################################
die() {
  err "$@"
  exit 1
}

#######################################
# Prints how to call the calling script and exits with status 2, as a bad
# flag would.
# Arguments:
#   The script's arguments, e.g. "<collection-id>".
# Outputs:
#   Writes "usage: <script> <arguments>" to STDERR.
#######################################
usage() {
  echo "usage: $(basename "$0") $*" >&2
  exit 2
}

#######################################
# Exits with status 127 when a command isn't installed.
# Arguments:
#   The command name or path.
#   A hint on how to install it.
# Outputs:
#   Writes the error and the hint to STDERR when the command is missing.
#######################################
require_cmd() {
  local cmd="$1"
  local hint="$2"

  if ! command -v "${cmd}" > /dev/null 2>&1; then
    err "${cmd} not found. ${hint}"
    exit 127
  fi
}

#######################################
# Checks that linkwarden and jq are installed.
# Globals:
#   LW
# Outputs:
#   Writes an error and an install hint to STDERR when one is missing.
#######################################
check_deps() {
  require_cmd "${LW}" \
    "Install linkwarden, or set LINKWARDEN_BIN to a local build."
  require_cmd jq "Install it with: brew install jq"
}

#######################################
# Exits with status 2, as a usage error, unless a value is a positive whole
# number.
# Arguments:
#   The value's name, used in the error, e.g. "collection-id".
#   The value.
# Outputs:
#   Writes the error to STDERR when the value isn't a positive number.
#######################################
require_id() {
  local name="$1"
  local value="$2"

  if [[ ! "${value}" =~ ^[1-9][0-9]*$ ]]; then
    err "${name} must be a positive number, got '${value}'"
    exit 2
  fi
}
