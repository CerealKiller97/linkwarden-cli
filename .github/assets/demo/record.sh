#!/usr/bin/env bash
#
# Records the README's demo GIFs with VHS, end to end:
#
#   1. starts a throwaway Linkwarden in Docker (compose.yml, fresh volumes)
#   2. creates a demo user and API token, and seeds it with seed.sh
#   3. builds the CLI for Linux into bin/, and a VHS image with curl and jq
#   4. records each tape into .github/assets/
#   5. removes the Linkwarden containers and volumes again
#
# Usage: .github/assets/demo/record.sh [tape...]
#   .github/assets/demo/record.sh                    # every demo tape
#   .github/assets/demo/record.sh collections links  # only these
#   KEEP=1 .github/assets/demo/record.sh links       # leave Linkwarden running
#
# Environment: DEMO_PORT (3200), LINKWARDEN_VERSION (latest), KEEP (unset).
# Needs docker, go, curl, jq and openssl. Never touches your own account: the
# CLI only ever talks to the throwaway instance.

set -euo pipefail

DEMO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${DEMO_DIR}/../../.." && pwd)"
readonly DEMO_DIR ROOT_DIR
readonly ALL_TAPES=(interactive collections links import scripting)
readonly VHS_IMAGE="linkwarden-cli-vhs"
readonly PORT="${DEMO_PORT:-3200}"
readonly BASE_URL="http://localhost:${PORT}"

export DEMO_PORT="${PORT}"
POSTGRES_PASSWORD="$(openssl rand -hex 16)"
MEILI_MASTER_KEY="$(openssl rand -hex 16)"
NEXTAUTH_SECRET="$(openssl rand -hex 32)"
export POSTGRES_PASSWORD MEILI_MASTER_KEY NEXTAUTH_SECRET

#######################################
# Prints a step heading to STDERR.
# Arguments:
#   Message.
#######################################
step() {
  printf '\033[1;35m==>\033[0m %s\n' "$*" >&2
}

#######################################
# Exits with an error message.
# Arguments:
#   Message.
#######################################
die() {
  echo "error: $*" >&2
  exit 1
}

#######################################
# Runs docker compose on the demo stack.
# Arguments:
#   Compose arguments.
#######################################
compose() {
  docker compose -f "${DEMO_DIR}/compose.yml" "$@"
}

#######################################
# Removes the demo stack, unless KEEP is set.
#######################################
cleanup() {
  if [[ -n "${KEEP:-}" ]]; then
    step "Leaving Linkwarden running at ${BASE_URL} (stop it with: docker compose -f ${DEMO_DIR}/compose.yml down -v)"
    return
  fi
  step "Removing the demo Linkwarden"
  compose down -v >/dev/null 2>&1 || true
}

#######################################
# Checks the tools and the tape names.
# Arguments:
#   Tape names.
#######################################
check() {
  local bin
  for bin in docker go curl jq openssl; do
    command -v "${bin}" >/dev/null || die "${bin} is required"
  done
  docker info >/dev/null 2>&1 || die "the Docker daemon isn't running"

  local tape
  for tape in "$@"; do
    [[ -f "${DEMO_DIR}/${tape}.tape" ]] || die "no such tape: ${tape} (have: ${ALL_TAPES[*]})"
  done
}

#######################################
# Starts a fresh Linkwarden and waits until it answers.
#######################################
start_linkwarden() {
  step "Starting a throwaway Linkwarden at ${BASE_URL}"
  compose down -v >/dev/null 2>&1 || true
  compose up -d --quiet-pull >/dev/null

  local _
  for _ in $(seq 1 120); do
    if curl -fsS -o /dev/null "${BASE_URL}/api/v1/logins" 2>/dev/null; then
      return
    fi
    sleep 2
  done
  die "Linkwarden didn't start in 4 minutes; see: docker compose -f ${DEMO_DIR}/compose.yml logs linkwarden"
}

#######################################
# Creates the demo user and an API token for it.
# Outputs:
#   The token on STDOUT.
#######################################
create_token() {
  # Registration can fail if a retry already created the user; signing in below
  # is what has to work.
  curl -sS -X POST "${BASE_URL}/api/v1/users" -H 'Content-Type: application/json' \
    -d '{"name":"Demo","username":"demo","password":"demopassword"}' >/dev/null || true

  local session
  session="$(curl -fsS -X POST "${BASE_URL}/api/v1/session" -H 'Content-Type: application/json' \
    -d '{"username":"demo","password":"demopassword","sessionName":"vhs-session"}' | jq -r '.response.token')"

  # The session is stored as a token named after sessionName, so this one needs
  # another name. expires: 4 = never
  curl -fsS -X POST "${BASE_URL}/api/v1/tokens" -H 'Content-Type: application/json' \
    -H "Authorization: Bearer ${session}" -d '{"name":"vhs-recording","expires":4}' | jq -r '.response.secretKey'
}

#######################################
# Builds the CLI for the host (to seed) and for Docker's Linux (to record).
# Globals:
#   HOST_BIN
#######################################
build_cli() {
  step "Building the CLI"
  local arch
  case "$(docker info --format '{{.Architecture}}')" in
    x86_64 | amd64) arch=amd64 ;;
    aarch64 | arm64) arch=arm64 ;;
    *) die "unsupported Docker architecture" ;;
  esac

  mkdir -p "${ROOT_DIR}/bin"
  HOST_BIN="${ROOT_DIR}/bin/linkwarden-host"
  (cd "${ROOT_DIR}" && go build -o "${HOST_BIN}" ./cmd/linkwarden)
  (cd "${ROOT_DIR}" && CGO_ENABLED=0 GOOS=linux GOARCH="${arch}" \
    go build -ldflags "-s -w" -o bin/linkwarden-demo ./cmd/linkwarden)
}

#######################################
# Builds the VHS image: the official one plus curl and jq.
#######################################
build_vhs_image() {
  step "Building the ${VHS_IMAGE} image"
  printf 'FROM ghcr.io/charmbracelet/vhs\nRUN apt-get update && apt-get install -y --no-install-recommends curl jq && rm -rf /var/lib/apt/lists/*\n' \
    | docker build -q -t "${VHS_IMAGE}" - >/dev/null
}

#######################################
# Records one tape.
# Arguments:
#   Tape name.
#######################################
record() {
  step "Recording $1.gif"
  docker run --rm --shm-size=1g --add-host=host.docker.internal:host-gateway \
    -v "${ROOT_DIR}:/vhs" \
    -e LINKWARDEN_URL="http://host.docker.internal:${PORT}" -e LINKWARDEN_TOKEN \
    "${VHS_IMAGE}" ".github/assets/demo/$1.tape" >/dev/null
}

main() {
  local tapes=("$@")
  if [[ ${#tapes[@]} -eq 0 ]]; then
    tapes=("${ALL_TAPES[@]}")
  fi
  check "${tapes[@]}"

  build_cli
  build_vhs_image

  trap cleanup EXIT
  start_linkwarden

  step "Seeding demo data"
  LINKWARDEN_TOKEN="$(create_token)"
  export LINKWARDEN_TOKEN
  LINKWARDEN_URL="${BASE_URL}" XDG_STATE_HOME="$(mktemp -d)" LINKWARDEN_BIN="${HOST_BIN}" \
    "${DEMO_DIR}/seed.sh" >&2

  local tape
  for tape in "${tapes[@]}"; do
    record "${tape}"
  done
  step "Done: ${tapes[*]/%/.gif} in .github/assets/"
}

main "$@"
