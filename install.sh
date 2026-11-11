#!/usr/bin/env bash
# shellcheck enable=require-variable-braces,quote-safe-variables,check-unassigned-uppercase
#
# 🔗 linkwarden-cli installer. Compatible with bash 3.2 (macOS's /bin/bash).
#
# Finds the release for this OS and CPU on GitHub, checks it against the
# release's checksums.txt and installs the linkwarden binary into
# /usr/local/bin (asking sudo when needed), or ~/.local/bin if you pick it,
# plus tab completions for your shell. If the folder isn't on your PATH, it
# offers to add it.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/CerealKiller97/linkwarden-cli/main/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/CerealKiller97/linkwarden-cli/main/install.sh | bash -s -- --yes
#
# Run with --help for the options.
#
# Everything happens inside main, called on the last line, so a download cut
# off halfway runs nothing.

set -euo pipefail

readonly REPO="CerealKiller97/linkwarden-cli"
readonly BIN="linkwarden"
readonly STEPS=6
readonly SYSTEM_BIN_DIR="/usr/local/bin"
readonly SYSTEM_SHARE_DIR="/usr/local/share"

# Options, set by parse_args.
OPT_VERSION="${LINKWARDEN_VERSION:-}"
OPT_DIR="${LINKWARDEN_INSTALL_DIR:-}"
RELEASES_URL="${LINKWARDEN_RELEASES_URL:-https://github.com/${REPO}/releases}"
ASSUME_YES=0
PLAIN=0
QUIET=0
FORCE=0
NO_COMPLETIONS=0
NO_MODIFY_PATH=0
UNINSTALL=0

# State, filled in as the install goes.
FANCY=0
TRUECOLOR=0
INTERACTIVE=0
TMP_DIR=""
OS=""
ARCH=""
VERSION=""
ASSET=""
INSTALL_DIR=""
SUDO=""
SHELL_NAME=""
EXISTING_PATH=""
EXISTING_VERSION=""
GRAD=""
REPLY_TEXT=""

C_RESET=""
C_BOLD=""
C_DIM=""
C_ITALIC=""
C_RED=""
C_GREEN=""
C_YELLOW=""
C_CYAN=""
C_WHITE=""

SPINNER=("⠋" "⠙" "⠹" "⠸" "⠼" "⠴" "⠦" "⠧" "⠇" "⠏")

# "linkwarden" in box-drawing letters, one array entry per letter and row.
LOGO_TOP=("╻  " "╻" "┏┓╻" "╻┏ " "╻ ╻" "┏━┓" "┏━┓" "╺┳┓" "┏━╸" "┏┓╻")
LOGO_MID=("┃  " "┃" "┃┗┫" "┣┻┓" "┃╻┃" "┣━┫" "┣┳┛" " ┃┃" "┣╸ " "┃┗┫")
LOGO_BOT=("┗━╸" "╹" "╹ ╹" "╹ ╹" "┗┻┛" "╹ ╹" "╹┗╸" "╺┻┛" "┗━╸" "╹ ╹")

CONFETTI=("*" "+" "✦" "✧" "•" "·" "°" "✶" "⋆")

#######################################
# Prints the usage text.
# Outputs:
#   Writes the help to STDOUT.
#######################################
usage() {
  cat <<EOF
🔗 linkwarden-cli installer

Usage:
  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash
  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- [options]
  bash install.sh [options]

Options:
  -v, --version <x.y.z>  Install this version instead of the latest.
  -d, --dir <path>       Install into <path>. Default: /usr/local/bin. Without a terminal,
                         falls back to \$XDG_BIN_HOME or ~/.local/bin when sudo would
                         need a password.
  -y, --yes              Don't ask anything; take the defaults. Also on when there's no
                         terminal or CI is set. Never edits your shell rc file.
  -f, --force            Reinstall even when that version is already installed.
      --no-completions   Don't install tab completions (installed by default).
      --no-modify-path   Don't offer to add the install folder to your PATH.
      --plain            No colours or animations. Also on with NO_COLOR or TERM=dumb.
  -q, --quiet            Only print errors and the final result.
      --uninstall        Remove linkwarden and the completions this script installs.
  -h, --help             Show this help.

Environment:
  LINKWARDEN_VERSION       Same as --version.
  LINKWARDEN_INSTALL_DIR   Same as --dir.
  LINKWARDEN_RELEASES_URL  Base URL of the releases (default: ${RELEASES_URL}).
EOF
}

#######################################
# Prints an error and exits with status 1.
# Arguments:
#   The message.
# Outputs:
#   Writes the message to STDERR.
#######################################
die() {
  if (( FANCY )); then
    printf '\r\033[K'
  fi
  printf '\n  %s✖ %s%s\n\n' "${C_RED}${C_BOLD}" "$*" "${C_RESET}" >&2
  exit 1
}

#######################################
# Prints a warning.
# Arguments:
#   The message.
# Outputs:
#   Writes the message to STDERR.
#######################################
warn() {
  printf '  %s⚠%s  %s\n' "${C_YELLOW}${C_BOLD}" "${C_RESET}" "$*" >&2
}

#######################################
# Prints an info line, unless --quiet.
# Arguments:
#   The message.
#######################################
info() {
  if (( ! QUIET )); then
    printf '  %s\n' "$*"
  fi
}

#######################################
# Prints a success line, unless --quiet.
# Arguments:
#   The message.
#######################################
ok() {
  if (( ! QUIET )); then
    printf '  %s✔%s %s\n' "${C_GREEN}${C_BOLD}" "${C_RESET}" "$*"
  fi
}

#######################################
# Prints a dim, italic aside, unless --quiet.
# Arguments:
#   The message.
#######################################
quip() {
  if (( ! QUIET )); then
    printf '  %s%s%s\n' "${C_DIM}${C_ITALIC}" "$*" "${C_RESET}"
  fi
}

#######################################
# Sleeps only when animating.
# Arguments:
#   Seconds, may be fractional.
#######################################
pause() {
  if (( FANCY )); then
    sleep "$1"
  fi
}

#######################################
# Replaces a leading $HOME with ~ for display.
# Arguments:
#   A path.
# Outputs:
#   Writes the shortened path to STDOUT.
#######################################
tildify() {
  case "$1" in
    "${HOME}") echo "~" ;;
    "${HOME}"/*) echo "~${1#"${HOME}"}" ;;
    *) echo "$1" ;;
  esac
}

#######################################
# Parses the command-line options into the OPT_* and flag globals.
# Arguments:
#   The script's arguments.
#######################################
parse_args() {
  while (( $# > 0 )); do
    case "$1" in
      -v | --version)
        [[ $# -ge 2 ]] || die "$1 needs a version, e.g. $1 1.0.0"
        OPT_VERSION="$2"
        shift
        ;;
      --version=*) OPT_VERSION="${1#*=}" ;;
      -d | --dir)
        [[ $# -ge 2 ]] || die "$1 needs a folder"
        OPT_DIR="$2"
        shift
        ;;
      --dir=*) OPT_DIR="${1#*=}" ;;
      -y | --yes) ASSUME_YES=1 ;;
      -f | --force) FORCE=1 ;;
      --no-completions) NO_COMPLETIONS=1 ;;
      --no-modify-path) NO_MODIFY_PATH=1 ;;
      --plain) PLAIN=1 ;;
      -q | --quiet) QUIET=1 ;;
      --uninstall) UNINSTALL=1 ;;
      -h | --help)
        usage
        exit 0
        ;;
      *)
        usage >&2
        die "unknown option: $1"
        ;;
    esac
    shift
  done
  OPT_VERSION="${OPT_VERSION#v}"
}

#######################################
# Decides whether to use colours, animations and prompts, and sets the
# colour variables.
#######################################
setup_terminal() {
  if [[ -t 1 && -z "${NO_COLOR:-}" && "${TERM:-dumb}" != "dumb" ]] \
    && (( ! PLAIN && ! QUIET )); then
    FANCY=1
  fi
  case "${COLORTERM:-}" in
    truecolor | 24bit) TRUECOLOR=1 ;;
    *) ;;
  esac
  if (( ! ASSUME_YES )) && [[ -z "${CI:-}" && -t 1 ]] \
    && { : </dev/tty; } 2>/dev/null; then
    INTERACTIVE=1
  fi
  if (( FANCY )); then
    C_RESET=$'\033[0m'
    C_BOLD=$'\033[1m'
    C_DIM=$'\033[2m'
    C_ITALIC=$'\033[3m'
    C_RED=$'\033[31m'
    C_GREEN=$'\033[32m'
    C_YELLOW=$'\033[33m'
    C_CYAN=$'\033[36m'
    C_WHITE=$'\033[97m'
  fi
}

#######################################
# Cleans up on exit: removes the temp folder and shows the cursor again.
#######################################
cleanup() {
  if [[ -n "${TMP_DIR}" ]]; then
    rm -rf "${TMP_DIR}"
  fi
  if (( FANCY )); then
    printf '\033[?25h'
  fi
}

#######################################
# Handles Ctrl-C: stops background downloads and exits with status 130.
#######################################
interrupted() {
  local pids
  pids="$(jobs -p)"
  if [[ -n "${pids}" ]]; then
    # shellcheck disable=SC2086 # One PID per word.
    kill ${pids} 2>/dev/null || true
  fi
  if (( FANCY )); then
    printf '\r\033[K'
  fi
  printf '\n  %sInstall cancelled. Nothing was changed after this point. 👋%s\n\n' \
    "${C_YELLOW}" "${C_RESET}" >&2
  exit 130
}

#######################################
# Sets GRAD to the escape code for a point on the sky → indigo → fuchsia
# gradient. Uses printf -v, not a subshell, so it's cheap to call per letter.
# Arguments:
#   Position, from 0.
#   Number of positions.
#######################################
grad_at() {
  local pos="$1" total="$2" t r g b
  if (( total < 2 )); then
    total=2
  fi
  t=$(( pos * 1000 / (total - 1) ))
  if (( t <= 500 )); then
    r=$(( 56 + (129 - 56) * t / 500 ))
    g=$(( 189 + (140 - 189) * t / 500 ))
    b=$(( 248 + (248 - 248) * t / 500 ))
  else
    t=$(( t - 500 ))
    r=$(( 129 + (232 - 129) * t / 500 ))
    g=$(( 140 + (121 - 140) * t / 500 ))
    b=$(( 248 + (249 - 248) * t / 500 ))
  fi
  if (( TRUECOLOR )); then
    printf -v GRAD '\033[38;2;%d;%d;%dm' "${r}" "${g}" "${b}"
  else
    printf -v GRAD '\033[38;5;%dm' \
      $(( 16 + 36 * ((r * 5 + 127) / 255) + 6 * ((g * 5 + 127) / 255) + (b * 5 + 127) / 255 ))
  fi
}

#######################################
# Prints one row of the logo, letter by letter along the gradient. The letter
# at the highlight index is drawn bright white, for the shimmer.
# Arguments:
#   Row: top, mid or bot.
#   Highlighted letter index, or -1.
#######################################
logo_row() {
  local row="$1" hl="$2" i n=10 out="  " glyph
  for (( i = 0; i < n; i++ )); do
    case "${row}" in
      top) glyph="${LOGO_TOP[i]}" ;;
      mid) glyph="${LOGO_MID[i]}" ;;
      *) glyph="${LOGO_BOT[i]}" ;;
    esac
    if (( i == hl )); then
      out+="${C_BOLD}${C_WHITE}${glyph}${C_RESET}"
    else
      grad_at "${i}" "${n}"
      out+="${C_BOLD}${GRAD}${glyph}${C_RESET}"
    fi
  done
  if [[ "${row}" == "mid" ]]; then
    out+="  ${C_DIM}cli${C_RESET}"
  fi
  printf '%s\033[K\n' "${out}"
}

#######################################
# Prints text one character at a time, like it's being typed.
# Arguments:
#   The text. Keep it ASCII: bash 3.2 may split multibyte characters.
#######################################
typewrite() {
  local text="$1" i
  if (( ! FANCY )); then
    printf '%s\n' "${text}"
    return
  fi
  for (( i = 0; i < ${#text}; i++ )); do
    printf '%s' "${text:i:1}"
    sleep 0.012
  done
  printf '\n'
}

#######################################
# The opening: the logo drops in row by row, a shimmer sweeps across it, and
# the tagline types itself out.
#######################################
intro() {
  local hl
  if (( QUIET )); then
    return
  fi
  printf '\n'
  if (( ! FANCY )); then
    printf '  linkwarden cli installer\n\n'
    return
  fi
  printf '\033[?25l'
  logo_row top -1
  sleep 0.08
  logo_row mid -1
  sleep 0.08
  logo_row bot -1
  for (( hl = 0; hl <= 10; hl++ )); do
    printf '\033[3A'
    logo_row top "${hl}"
    logo_row mid "${hl}"
    logo_row bot "${hl}"
    sleep 0.05
  done
  printf '\n  🔖 '
  typewrite "Your bookmarks, from the comfort of your terminal."
  printf '\n'
}

#######################################
# Prints a step heading such as "◆ 2/6  Finding the latest release".
# Arguments:
#   Step number.
#   Title.
#######################################
step() {
  if (( QUIET )); then
    return
  fi
  grad_at $(( $1 - 1 )) "${STEPS}"
  if (( FANCY )); then
    printf '\n%s%s◆ %s/%s%s  %s%s%s\n' "${C_BOLD}" "${GRAD}" "$1" "${STEPS}" \
      "${C_RESET}" "${C_BOLD}" "$2" "${C_RESET}"
  else
    printf '\n[%s/%s] %s\n' "$1" "${STEPS}" "$2"
  fi
}

#######################################
# Runs a command in the background while a spinner turns. The command's
# STDOUT lands in ${TMP_DIR}/spin.out, its STDERR in ${TMP_DIR}/spin.err.
# Arguments:
#   Label shown next to the spinner.
#   The command and its arguments.
# Returns:
#   The command's exit status.
#######################################
spin() {
  local label="$1" pid i=0 status=0
  shift
  "$@" >"${TMP_DIR}/spin.out" 2>"${TMP_DIR}/spin.err" &
  pid=$!
  if (( FANCY )); then
    # Turn at least a few frames, so quick steps don't just flash by.
    while kill -0 "${pid}" 2>/dev/null || (( i < 8 )); do
      printf '\r  %s%s%s %s' "${C_CYAN}" "${SPINNER[i % ${#SPINNER[@]}]}" "${C_RESET}" "${label}"
      i=$(( i + 1 ))
      sleep 0.07
    done
    printf '\r\033[K'
  fi
  wait "${pid}" || status=$?
  return "${status}"
}

#######################################
# Asks a yes/no question on the terminal.
# Arguments:
#   The question.
#   The default answer: y or n.
# Returns:
#   0 for yes, 1 for no. Without a terminal, the default.
#######################################
confirm() {
  local question="$1" default="$2" hint answer
  if (( ! INTERACTIVE )); then
    [[ "${default}" == "y" ]]
    return
  fi
  if [[ "${default}" == "y" ]]; then
    hint="Y/n"
  else
    hint="y/N"
  fi
  printf '\033[?25h'
  printf '  %s?%s %s %s(%s)%s ' "${C_CYAN}${C_BOLD}" "${C_RESET}" "${question}" \
    "${C_DIM}" "${hint}" "${C_RESET}"
  read -r answer </dev/tty || answer=""
  if (( FANCY )); then
    printf '\033[?25l'
  fi
  answer="${answer:-${default}}"
  case "${answer}" in
    [Yy]*) return 0 ;;
    *) return 1 ;;
  esac
}

#######################################
# Reads a line of text from the terminal into REPLY_TEXT.
# Arguments:
#   The prompt.
#######################################
ask() {
  printf '\033[?25h'
  printf '  %s?%s %s ' "${C_CYAN}${C_BOLD}" "${C_RESET}" "$1"
  read -r REPLY_TEXT </dev/tty || REPLY_TEXT=""
  if (( FANCY )); then
    printf '\033[?25l'
  fi
}

#######################################
# Exits with an error when a command the installer needs is missing.
#######################################
require_tools() {
  local cmd
  for cmd in curl tar uname mktemp; do
    if ! command -v "${cmd}" >/dev/null 2>&1; then
      die "this installer needs '${cmd}', which isn't installed"
    fi
  done
  if ! command -v sha256sum >/dev/null 2>&1 && ! command -v shasum >/dev/null 2>&1 \
    && ! command -v openssl >/dev/null 2>&1; then
    die "this installer needs sha256sum, shasum or openssl to check the download"
  fi
}

#######################################
# Sets OS and ARCH to GoReleaser's names for this machine, and SHELL_NAME.
#######################################
detect_platform() {
  local os arch
  os="$(uname -s)"
  arch="$(uname -m)"
  case "${os}" in
    Linux) OS="linux" ;;
    Darwin) OS="darwin" ;;
    MINGW* | MSYS* | CYGWIN* | Windows_NT)
      die "On Windows, run this in PowerShell instead:  irm https://raw.githubusercontent.com/${REPO}/main/install.ps1 | iex"
      ;;
    *) die "no prebuilt linkwarden for ${os}. Try: go install github.com/${REPO}/cmd/linkwarden@latest" ;;
  esac
  case "${arch}" in
    x86_64 | amd64) ARCH="amd64" ;;
    arm64 | aarch64 | armv8*) ARCH="arm64" ;;
    *) die "no prebuilt linkwarden for ${arch} CPUs. Try: go install github.com/${REPO}/cmd/linkwarden@latest" ;;
  esac
  # Rosetta 2 reports x86_64 to a translated shell; the native build is faster.
  if [[ "${OS}" == "darwin" && "${ARCH}" == "amd64" ]] \
    && [[ "$(sysctl -n sysctl.proc_translated 2>/dev/null || echo 0)" == "1" ]]; then
    ARCH="arm64"
  fi
  SHELL_NAME="$(basename "${SHELL:-sh}")"
}

#######################################
# Prints what detect_platform found, with a little commentary.
#######################################
show_platform() {
  local os_label cpu_label os_icon
  if [[ "${OS}" == "darwin" ]]; then
    os_icon="🍎"
    os_label="macOS"
  else
    os_icon="🐧"
    os_label="Linux"
  fi
  if [[ "${ARCH}" == "arm64" ]]; then
    if [[ "${OS}" == "darwin" ]]; then
      cpu_label="Apple Silicon"
    else
      cpu_label="ARM64"
    fi
  else
    cpu_label="x86_64"
  fi
  pause 0.15
  info "${os_icon} ${C_DIM}OS   ${C_RESET} ${C_BOLD}${os_label}${C_RESET} ${C_DIM}(${OS})${C_RESET}"
  pause 0.15
  info "🧠 ${C_DIM}CPU  ${C_RESET} ${C_BOLD}${cpu_label}${C_RESET} ${C_DIM}(${ARCH})${C_RESET}"
  pause 0.15
  info "🐚 ${C_DIM}Shell${C_RESET} ${C_BOLD}${SHELL_NAME}${C_RESET}"
  pause 0.15
  case "${OS}/${ARCH}" in
    darwin/arm64) quip "Apple Silicon. Fancy. 🚀" ;;
    darwin/amd64) quip "An Intel Mac. A classic. 🕰️" ;;
    linux/arm64) quip "A penguin with an ARM. Respect. 💪" ;;
    *) quip "Penguin spotted. Good taste. 🐧" ;;
  esac
}

#######################################
# Prints the latest release's version (without "v"), following the
# releases/latest redirect, so no API token or rate limit is involved.
# Outputs:
#   Writes the version to STDOUT.
#######################################
fetch_latest_version() {
  local url tag
  url="$(curl -fsSLI -o /dev/null -w '%{url_effective}' "${RELEASES_URL}/latest")"
  tag="${url##*/}"
  if [[ "${tag}" != v[0-9]* ]]; then
    echo "no release found at ${RELEASES_URL}/latest" >&2
    return 1
  fi
  echo "${tag#v}"
}

#######################################
# Sets VERSION and ASSET, from --version or the latest release.
#######################################
resolve_version() {
  if [[ -n "${OPT_VERSION}" ]]; then
    VERSION="${OPT_VERSION}"
    ok "Version ${C_BOLD}${VERSION}${C_RESET}, as requested"
  else
    if ! spin "Asking GitHub what's new…" fetch_latest_version; then
      die "couldn't find the latest release: $(cat "${TMP_DIR}/spin.err")"
    fi
    VERSION="$(cat "${TMP_DIR}/spin.out")"
    ok "Latest release is ${C_BOLD}v${VERSION}${C_RESET}"
    quip "Fresh out of the oven. 🥐"
  fi
  ASSET="${BIN}_${VERSION}_${OS}_${ARCH}.tar.gz"
}

#######################################
# Looks for an installed linkwarden and stops early when it's already this
# version, unless --force or the user asks to reinstall.
#######################################
check_existing() {
  EXISTING_PATH="$(command -v "${BIN}" 2>/dev/null || true)"
  if [[ -z "${EXISTING_PATH}" ]]; then
    return
  fi
  EXISTING_VERSION="$("${EXISTING_PATH}" --version 2>/dev/null | awk '{print $NF}' || true)"
  EXISTING_VERSION="${EXISTING_VERSION#v}"
  if [[ "${EXISTING_VERSION}" != "${VERSION}" ]]; then
    info "⬆️  Upgrade: ${C_DIM}${EXISTING_VERSION:-unknown}${C_RESET} → ${C_BOLD}${C_GREEN}${VERSION}${C_RESET}" \
      "${C_DIM}($(tildify "${EXISTING_PATH}"))${C_RESET}"
    return
  fi
  ok "linkwarden ${VERSION} is already installed at $(tildify "${EXISTING_PATH}")"
  if (( FORCE )); then
    return
  fi
  if confirm "Reinstall it anyway?" n; then
    return
  fi
  printf '\n  %sAll good, nothing to do. Happy bookmarking! 🔖%s\n\n' "${C_GREEN}" "${C_RESET}"
  exit 0
}

#######################################
# Expands a leading ~ in a path.
# Arguments:
#   A path.
# Outputs:
#   Writes the expanded path to STDOUT.
#######################################
expand_path() {
  # shellcheck disable=SC2088 # Matching a literal ~ typed at the prompt is the point.
  case "$1" in
    "~") echo "${HOME}" ;;
    "~/"*) echo "${HOME}/${1#"~/"}" ;;
    *) echo "$1" ;;
  esac
}

#######################################
# Sets SUDO to "sudo" when the folder (or its nearest existing parent) isn't
# writable by the current user.
# Arguments:
#   A folder.
#   Optional "soft": set SUDO to "missing" instead of exiting when sudo
#   isn't installed.
#######################################
pick_sudo() {
  local dir="$1" mode="${2:-}"
  SUDO=""
  while [[ ! -d "${dir}" && "${dir}" != "/" && -n "${dir}" ]]; do
    dir="$(dirname "${dir}")"
  done
  if [[ ! -w "${dir}" ]]; then
    if ! command -v sudo >/dev/null 2>&1; then
      if [[ "${mode}" == "soft" ]]; then
        SUDO="missing"
        return
      fi
      die "can't write to ${dir} and sudo isn't available. Pick another folder with --dir."
    fi
    SUDO="sudo"
  fi
}

#######################################
# Prints the per-user install folder: $XDG_BIN_HOME or ~/.local/bin.
# Outputs:
#   Writes the folder to STDOUT.
#######################################
user_bin_dir() {
  echo "${XDG_BIN_HOME:-${HOME}/.local/bin}"
}

#######################################
# Sets INSTALL_DIR: --dir, or a choice from a small menu, or the default,
# /usr/local/bin. Without a terminal, when /usr/local/bin needs sudo and sudo
# would ask for a password, falls back to the per-user folder.
#######################################
choose_dir() {
  local user_dir choice
  user_dir="$(user_bin_dir)"
  if [[ -n "${OPT_DIR}" ]]; then
    INSTALL_DIR="$(expand_path "${OPT_DIR}")"
  elif (( INTERACTIVE )); then
    printf '  %sWhere should linkwarden live?%s\n' "${C_BOLD}" "${C_RESET}"
    printf '    %s1)%s %-22s %s← everyone on this machine, may ask for sudo (default)%s\n' \
      "${C_CYAN}" "${C_RESET}" "${SYSTEM_BIN_DIR}" "${C_DIM}" "${C_RESET}"
    printf '    %s2)%s %-22s %s← just for you, no sudo%s\n' "${C_CYAN}" "${C_RESET}" \
      "$(tildify "${user_dir}")" "${C_DIM}" "${C_RESET}"
    printf '    %s3)%s %-22s\n' "${C_CYAN}" "${C_RESET}" "somewhere else…"
    ask "Pick one [1-3]:"
    choice="${REPLY_TEXT:-1}"
    case "${choice}" in
      1) INSTALL_DIR="${SYSTEM_BIN_DIR}" ;;
      2) INSTALL_DIR="${user_dir}" ;;
      3)
        ask "Folder:"
        [[ -n "${REPLY_TEXT}" ]] || die "no folder given"
        INSTALL_DIR="$(expand_path "${REPLY_TEXT}")"
        ;;
      *) die "'${choice}' isn't one of 1, 2 or 3" ;;
    esac
  else
    INSTALL_DIR="${SYSTEM_BIN_DIR}"
    pick_sudo "${INSTALL_DIR}" soft
    if [[ "${SUDO}" == "missing" ]] \
      || { [[ "${SUDO}" == "sudo" ]] && ! sudo -n true 2>/dev/null; }; then
      warn "Can't write to ${SYSTEM_BIN_DIR} without sudo and a terminal; using $(tildify "${user_dir}") instead"
      INSTALL_DIR="${user_dir}"
    fi
  fi
  INSTALL_DIR="${INSTALL_DIR%/}"
  pick_sudo "${INSTALL_DIR}"
  ok "Installing to ${C_BOLD}$(tildify "${INSTALL_DIR}")${C_RESET}${SUDO:+ ${C_DIM}(with sudo)${C_RESET}}"
  quip "Cozy. 🏠"
}

#######################################
# Prints the size of the file at a URL, following redirects, or nothing.
# Arguments:
#   The URL.
# Outputs:
#   Writes the byte count to STDOUT.
#######################################
content_length() {
  curl -fsSLI "$1" 2>/dev/null | tr -d '\r' \
    | awk 'tolower($1) == "content-length:" { n = $2 } END { if (n > 0) print n }' || true
}

#######################################
# Formats a byte count as e.g. "4.2 MB".
# Arguments:
#   Bytes.
# Outputs:
#   Writes the size to STDOUT.
#######################################
human_size() {
  awk -v b="$1" 'BEGIN {
    if (b >= 1048576) printf "%.1f MB", b / 1048576
    else if (b >= 1024) printf "%.0f KB", b / 1024
    else printf "%d B", b
  }'
}

#######################################
# Draws the download progress bar, with a bookmark running along it.
# Arguments:
#   Bytes so far.
#   Total bytes.
#######################################
draw_bar() {
  local done_bytes="$1" total="$2" width=32 filled pct i bar=""
  filled=$(( done_bytes * width / total ))
  if (( filled > width )); then
    filled="${width}"
  fi
  pct=$(( done_bytes * 100 / total ))
  for (( i = 0; i < width; i++ )); do
    if (( i < filled )); then
      grad_at "${i}" "${width}"
      bar+="${GRAD}━"
    elif (( i == filled )); then
      bar+="${C_RESET}🔖"
      i=$(( i + 1 ))
    else
      bar+="${C_DIM}─"
    fi
  done
  printf '\r  %s%s %s%3d%%%s %s%s / %s%s\033[K' "${bar}" "${C_RESET}" "${C_BOLD}" "${pct}" \
    "${C_RESET}" "${C_DIM}" "$(human_size "${done_bytes}")" "$(human_size "${total}")" "${C_RESET}"
}

#######################################
# Downloads a URL to a file, with a progress bar when there's a terminal.
# Arguments:
#   The URL.
#   The destination file.
#######################################
download() {
  local url="$1" dest="$2" total pid size status=0
  total=""
  if (( FANCY )); then
    total="$(content_length "${url}")"
  fi
  if [[ -z "${total}" ]]; then
    curl -fsSL -o "${dest}" "${url}" || return
    return
  fi
  curl -fsSL -o "${dest}" "${url}" 2>"${TMP_DIR}/curl.err" &
  pid=$!
  while kill -0 "${pid}" 2>/dev/null; do
    size=0
    if [[ -f "${dest}" ]]; then
      size=$(( $(wc -c <"${dest}") ))
    fi
    draw_bar "${size}" "${total}"
    sleep 0.05
  done
  wait "${pid}" || status=$?
  if (( status != 0 )); then
    printf '\r\033[K'
    cat "${TMP_DIR}/curl.err" >&2
    return "${status}"
  fi
  draw_bar "${total}" "${total}"
  printf '\n'
}

#######################################
# Downloads the release archive and its checksums.txt into TMP_DIR.
#######################################
fetch_release() {
  local base="${RELEASES_URL}/download/v${VERSION}"
  info "📦 ${C_DIM}${ASSET}${C_RESET}"
  if ! download "${base}/${ASSET}" "${TMP_DIR}/${ASSET}"; then
    die "couldn't download ${base}/${ASSET}. Is v${VERSION} a real release? See ${RELEASES_URL}"
  fi
  if ! spin "Grabbing checksums.txt…" curl -fsSL -o "${TMP_DIR}/checksums.txt" \
    "${base}/checksums.txt"; then
    die "couldn't download ${base}/checksums.txt"
  fi
  ok "Downloaded $(human_size "$(( $(wc -c <"${TMP_DIR}/${ASSET}") ))")"
}

#######################################
# Prints a file's SHA-256 with whichever tool is installed.
# Arguments:
#   The file.
# Outputs:
#   Writes the hex digest to STDOUT.
#######################################
sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1}'
  else
    openssl dgst -sha256 -r "$1" | awk '{print $1}'
  fi
}

#######################################
# Checks the archive against checksums.txt and exits when they differ.
#######################################
verify() {
  local expected actual
  expected="$(awk -v f="${ASSET}" '$2 == f || $2 == "*" f { print $1 }' "${TMP_DIR}/checksums.txt")"
  [[ -n "${expected}" ]] || die "checksums.txt has no entry for ${ASSET}"
  spin "Crunching SHA-256…" sha256_of "${TMP_DIR}/${ASSET}" || die "couldn't hash ${ASSET}"
  actual="$(cat "${TMP_DIR}/spin.out")"
  if [[ "${actual}" != "${expected}" ]]; then
    printf '\n  %sexpected %s%s\n  %sgot      %s%s\n' "${C_DIM}" "${expected}" "${C_RESET}" \
      "${C_DIM}" "${actual}" "${C_RESET}" >&2
    die "checksum mismatch for ${ASSET}. The download is corrupt or was tampered with; nothing was installed."
  fi
  ok "SHA-256 matches ${C_DIM}${actual:0:12}…${C_RESET}"
  quip "Nobody's been messing with your bookmarks. 🔒"
}

#######################################
# Extracts the archive and installs the binary into INSTALL_DIR.
#######################################
install_binary() {
  local src_dir="${TMP_DIR}/extract" src
  mkdir -p "${src_dir}"
  tar -xzf "${TMP_DIR}/${ASSET}" -C "${src_dir}" || die "couldn't extract ${ASSET}"
  src="${src_dir}/${ASSET%.tar.gz}/${BIN}"
  [[ -f "${src}" ]] || die "${ASSET} doesn't contain ${BIN}"
  if [[ "${OS}" == "darwin" ]]; then
    # The binary isn't notarized; don't let Gatekeeper block it.
    xattr -d com.apple.quarantine "${src}" 2>/dev/null || true
  fi
  if [[ -n "${SUDO}" ]]; then
    info "🔑 Asking sudo for permission to write to ${INSTALL_DIR}"
    printf '\033[?25h'
  fi
  ${SUDO} mkdir -p "${INSTALL_DIR}" || die "couldn't create ${INSTALL_DIR}"
  ${SUDO} install -m 0755 "${src}" "${INSTALL_DIR}/${BIN}" \
    || die "couldn't write ${INSTALL_DIR}/${BIN}"
  ok "Installed ${C_BOLD}$(tildify "${INSTALL_DIR}/${BIN}")${C_RESET}"
}

#######################################
# Prints where this script puts the completion file for a shell. System-wide
# paths are under /usr/local/share, which bash-completion, zsh and fish all
# search by default; per-user paths follow the XDG base directories.
# Arguments:
#   bash, zsh or fish.
#   system or user.
# Outputs:
#   Writes the path to STDOUT, or nothing for other shells.
#######################################
completion_path() {
  local data="${XDG_DATA_HOME:-${HOME}/.local/share}"
  local config="${XDG_CONFIG_HOME:-${HOME}/.config}"
  if [[ "$2" == "system" ]]; then
    case "$1" in
      bash) echo "${SYSTEM_SHARE_DIR}/bash-completion/completions/${BIN}" ;;
      zsh) echo "${SYSTEM_SHARE_DIR}/zsh/site-functions/_${BIN}" ;;
      fish) echo "${SYSTEM_SHARE_DIR}/fish/vendor_completions.d/${BIN}.fish" ;;
      *) ;;
    esac
    return
  fi
  case "$1" in
    bash) echo "${data}/bash-completion/completions/${BIN}" ;;
    zsh) echo "${data}/zsh/site-functions/_${BIN}" ;;
    fish) echo "${config}/fish/completions/${BIN}.fish" ;;
    *) ;;
  esac
}

#######################################
# Prints whether completions for INSTALL_DIR go system-wide or per-user:
# per-user when it's inside $HOME, system otherwise.
# Outputs:
#   Writes "system" or "user" to STDOUT.
#######################################
completion_scope() {
  case "${INSTALL_DIR}/" in
    "${HOME}"/*) echo "user" ;;
    *) echo "system" ;;
  esac
}

#######################################
# Installs completions for the user's shell, next to the binary: under
# /usr/local/share for a system-wide install, or the XDG folders for a
# per-user one. Skipped with --no-completions.
#######################################
install_completions() {
  local dest src scope comp_sudo
  if (( NO_COMPLETIONS )); then
    return
  fi
  scope="$(completion_scope)"
  dest="$(completion_path "${SHELL_NAME}" "${scope}")"
  if [[ -z "${dest}" ]]; then
    return
  fi
  src="${TMP_DIR}/extract/${ASSET%.tar.gz}/completions/${BIN}.${SHELL_NAME}"
  if [[ ! -f "${src}" ]]; then
    # Older releases don't ship completions/; generate them like
    # scripts/completions.sh does. Any token will do for `completion`.
    src="${TMP_DIR}/${BIN}.${SHELL_NAME}"
    if ! LINKWARDEN_TOKEN=completions "${INSTALL_DIR}/${BIN}" completion "${SHELL_NAME}" \
      >"${src}" 2>/dev/null || [[ ! -s "${src}" ]]; then
      return
    fi
  fi
  pick_sudo "$(dirname "${dest}")" soft
  comp_sudo="${SUDO}"
  pick_sudo "${INSTALL_DIR}"
  if [[ "${comp_sudo}" == "missing" ]]; then
    warn "Skipped tab completions: can't write to $(dirname "${dest}")"
    return
  fi
  if ! ${comp_sudo} mkdir -p "$(dirname "${dest}")" \
    || ! ${comp_sudo} install -m 0644 "${src}" "${dest}"; then
    warn "Skipped tab completions: couldn't write ${dest}"
    return
  fi
  ok "Tab completions for ${SHELL_NAME} → $(tildify "${dest}")"
  if [[ "${SHELL_NAME}" == "zsh" && "${scope}" == "user" ]]; then
    info "   ${C_DIM}zsh loads them once this is in ~/.zshrc, before compinit:${C_RESET}"
    info "   ${C_CYAN}fpath=($(tildify "$(dirname "${dest}")") \$fpath)${C_RESET}"
  fi
}

#######################################
# Prints the rc file to add a PATH line to, for the user's shell.
# Outputs:
#   Writes the path to STDOUT.
#######################################
rc_file() {
  case "${SHELL_NAME}" in
    zsh) echo "${ZDOTDIR:-${HOME}}/.zshrc" ;;
    bash)
      if [[ "${OS}" == "darwin" ]]; then
        echo "${HOME}/.bash_profile"
      else
        echo "${HOME}/.bashrc"
      fi
      ;;
    fish) echo "${XDG_CONFIG_HOME:-${HOME}/.config}/fish/conf.d/${BIN}.fish" ;;
    *) echo "${HOME}/.profile" ;;
  esac
}

#######################################
# When INSTALL_DIR isn't on PATH, offers to add it to the shell's rc file.
# Never edits anything without a yes on the terminal.
#######################################
ensure_path() {
  local rc line shown
  case ":${PATH}:" in
    *":${INSTALL_DIR}:"*) return ;;
    *) ;;
  esac
  shown="$(tildify "${INSTALL_DIR}")"
  rc="$(rc_file)"
  if [[ "${SHELL_NAME}" == "fish" ]]; then
    line="fish_add_path -g \"${INSTALL_DIR}\""
  else
    line="export PATH=\"${INSTALL_DIR}:\$PATH\""
  fi
  warn "${shown} isn't on your PATH yet"
  if (( NO_MODIFY_PATH )) || ! (( INTERACTIVE )) \
    || ! confirm "Add it to $(tildify "${rc}")?" y; then
    info "   Add it yourself with:  ${C_CYAN}echo '${line}' >> $(tildify "${rc}")${C_RESET}"
    return
  fi
  if [[ -f "${rc}" ]] && grep -qF "${line}" "${rc}"; then
    ok "$(tildify "${rc}") already has it; open a new terminal to pick it up"
    return
  fi
  mkdir -p "$(dirname "${rc}")"
  printf '\n# Added by the linkwarden-cli installer\n%s\n' "${line}" >>"${rc}"
  ok "Added to $(tildify "${rc}"). Open a new terminal (or source it) to use ${BIN}."
}

#######################################
# Confetti falling for a moment, then wiped away.
#######################################
confetti() {
  local rows=5 width=46 frame r c row
  local -a lines=()
  for (( r = 0; r < rows; r++ )); do
    lines[r]=""
    printf '\n'
  done
  for (( frame = 0; frame < 18; frame++ )); do
    row=""
    for (( c = 0; c < width; c++ )); do
      if (( RANDOM % 7 == 0 )); then
        grad_at $(( RANDOM % 10 )) 10
        row+="${GRAD}${CONFETTI[RANDOM % ${#CONFETTI[@]}]}"
      else
        row+=" "
      fi
    done
    # Shift every row down one and drop the new one in on top.
    for (( r = rows - 1; r > 0; r-- )); do
      lines[r]="${lines[r - 1]}"
    done
    lines[0]="${row}${C_RESET}"
    printf '\033[%dA' "${rows}"
    for (( r = 0; r < rows; r++ )); do
      printf '  %s\033[K\n' "${lines[r]}"
    done
    sleep 0.06
  done
  printf '\033[%dA' "${rows}"
  for (( r = 0; r < rows; r++ )); do
    printf '\033[K\n'
  done
  printf '\033[%dA' "${rows}"
}

#######################################
# The ending: confetti, proof the binary runs, and what to do next.
#######################################
finale() {
  local installed shadow
  installed="$("${INSTALL_DIR}/${BIN}" --version 2>/dev/null | awk '{print $NF}' || true)"
  if (( QUIET )); then
    printf 'linkwarden %s installed to %s\n' "${installed:-${VERSION}}" "${INSTALL_DIR}/${BIN}"
    return
  fi
  if (( FANCY )); then
    printf '\n'
    confetti
  fi
  printf '\n  🎉 %slinkwarden %s is ready!%s\n' "${C_BOLD}${C_GREEN}" "${installed:-${VERSION}}" \
    "${C_RESET}"
  shadow="$(command -v "${BIN}" 2>/dev/null || true)"
  if [[ -n "${shadow}" && "${shadow}" != "${INSTALL_DIR}/${BIN}" ]]; then
    warn "Another ${BIN} comes first on your PATH: $(tildify "${shadow}")"
    info "   Remove it or put $(tildify "${INSTALL_DIR}") earlier in PATH to use this one."
  fi
  printf '\n  %sNext up:%s\n' "${C_BOLD}" "${C_RESET}"
  printf '    %s$%s %slinkwarden init%s     %s# connect to your Linkwarden%s\n' "${C_DIM}" \
    "${C_RESET}" "${C_CYAN}" "${C_RESET}" "${C_DIM}" "${C_RESET}"
  printf '    %s$%s %slinkwarden --help%s   %s# look around%s\n' "${C_DIM}" "${C_RESET}" \
    "${C_CYAN}" "${C_RESET}" "${C_DIM}" "${C_RESET}"
  printf '\n  %sHappy bookmarking! 🔖%s\n\n' "${C_DIM}" "${C_RESET}"
}

#######################################
# Removes the binary from --dir, or from /usr/local/bin and the per-user
# folder, and any completions this script installs, system-wide or per-user.
#######################################
uninstall() {
  local dir target removed=0 shell scope file
  local -a dirs=()
  if [[ -n "${OPT_DIR}" ]]; then
    dirs=("$(expand_path "${OPT_DIR}")")
  else
    dirs=("${SYSTEM_BIN_DIR}" "$(user_bin_dir)")
  fi
  printf '\n'
  for dir in "${dirs[@]}"; do
    target="${dir%/}/${BIN}"
    if [[ ! -e "${target}" ]]; then
      info "No ${BIN} at $(tildify "${target}")"
      continue
    fi
    if ! confirm "Remove $(tildify "${target}")?" y; then
      info "Left it alone."
      continue
    fi
    pick_sudo "${dir}"
    ${SUDO} rm -f "${target}"
    ok "Removed $(tildify "${target}")"
    removed=1
  done
  for scope in system user; do
    for shell in bash zsh fish; do
      file="$(completion_path "${shell}" "${scope}")"
      if [[ -f "${file}" ]]; then
        pick_sudo "$(dirname "${file}")"
        ${SUDO} rm -f "${file}"
        ok "Removed $(tildify "${file}")"
        removed=1
      fi
    done
  done
  if [[ -n "$(command -v "${BIN}" 2>/dev/null || true)" ]]; then
    warn "${BIN} is still on your PATH at $(tildify "$(command -v "${BIN}")")"
  fi
  if (( removed )); then
    info "${C_DIM}Any PATH line the installer added to your shell rc file is still there.${C_RESET}"
    printf '\n  %sSorry to see you go. 👋%s\n\n' "${C_DIM}" "${C_RESET}"
  fi
}

#######################################
# Runs the install, step by step.
# Arguments:
#   The script's arguments.
#######################################
main() {
  parse_args "$@"
  setup_terminal
  TMP_DIR="$(mktemp -d 2>/dev/null || mktemp -d -t linkwarden)"
  trap cleanup EXIT
  trap interrupted INT TERM

  if (( UNINSTALL )); then
    uninstall
    return
  fi

  intro
  require_tools

  step 1 "Sniffing out your system"
  detect_platform
  show_platform

  step 2 "Finding the latest release"
  resolve_version
  check_existing

  step 3 "Picking a home"
  choose_dir

  step 4 "Downloading"
  fetch_release

  step 5 "Checking it's the real deal"
  verify

  step 6 "Moving in"
  install_binary
  install_completions
  ensure_path

  finale
}

main "$@"
