# Shared helpers for the example scripts. Source it, don't run it:
#   . "$(dirname "$0")/lib.sh"

# LW is the linkwarden binary: LINKWARDEN_BIN, or linkwarden from your PATH.
LW="${LINKWARDEN_BIN:-linkwarden}"

# die prints an error and exits 1.
die() {
  echo "error: $*" >&2
  exit 1
}

# usage prints how to call the script and exits 2, like a bad flag would.
usage() {
  echo "usage: $(basename "$0") $*" >&2
  exit 2
}

# require_cmd fails early with a hint when a command the script needs is missing.
require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: $1 not found. $2" >&2
    exit 127
  fi
}

# require_id exits 2, as a usage error, unless $2 is a positive whole number;
# $1 names it in the error.
require_id() {
  case "$2" in
    '' | *[!0-9]* | 0)
      echo "error: $1 must be a positive number, got '$2'" >&2
      exit 2
      ;;
  esac
}

require_cmd "$LW" "Install linkwarden, or set LINKWARDEN_BIN to a build (go build -o bin/linkwarden .)."
require_cmd jq "Install it with: brew install jq"
