#!/usr/bin/env bash
# 🚀 Release helper for linkwarden-cli. Compatible with bash 3.2 (macOS's /bin/bash).
#
# Usage:
#   scripts/release.sh <major|minor|patch> [-y|--yes]
#       Triggers the "release" GitHub Actions workflow for origin/main. Shows the next
#       version and a changelog preview first, then asks to confirm (-y skips it).
#       Requires the GitHub CLI (https://cli.github.com), authenticated with `gh auth login`.
#
#   The workflow drives the release through these subcommands. Run them locally to preview:
#
#   scripts/release.sh next-version <major|minor|patch>
#       Prints the next version (without "v") from the highest vX.Y.Z tag.
#       With no tags yet, the first release is 1.0.0.
#
#   scripts/release.sh changelog <version> [--ref REF] [--write]
#       Prints a Keep a Changelog (https://keepachangelog.com/en/1.1.0/) section built from
#       the Conventional Commits in <latest tag>..REF (REF defaults to HEAD). --write also
#       inserts it into CHANGELOG.md below "## [Unreleased]" and updates the compare links.
#
#   scripts/release.sh build <version> [out-dir]
#       Cross-compiles linux/darwin/windows × amd64/arm64 archives and a checksums.txt
#       into out-dir (default: dist/).
#
# Commit type → changelog section:
#   feat                                     → ✨ Added
#   refactor, perf, rename                   → ♻️ Changed
#   deprecate                                → ⚠️ Deprecated
#   remove                                   → 🗑️ Removed
#   fix                                      → 🐛 Fixed
#   security, or scope "security"            → 🔒 Security
#   docs, chore, ci, test, build, style      → omitted
#   and non-conventional commits
#
# A "!" after the type/scope (feat!: ..., fix(api)!: ...) or a "BREAKING CHANGE:" footer
# marks the entry with 💥 **BREAKING**.
set -euo pipefail

cd "$(dirname "$0")/.."

CHANGELOG="CHANGELOG.md"
WORKFLOW="release.yml"
BUILD_TARGETS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64"
VERSION_PKG="github.com/CerealKiller97/linkwarden-cli/cmd.version"

# ─── helpers ──────────────────────────────────────────────────────────────────

usage() {
	sed -n '2,/^set -euo/p' "$0" | sed -e '$d' -e 's/^# \{0,1\}//' >&2
	exit 1
}

die() {
	echo "❌ $*" >&2
	exit 1
}

# latest_tag prints the highest vX.Y.Z tag, or nothing if there is none.
latest_tag() {
	git tag --list 'v*' --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | head -n 1 || true
}

# repo_url prints the repository's https URL, e.g. https://github.com/owner/repo.
repo_url() {
	if [[ -n "${GITHUB_SERVER_URL:-}" && -n "${GITHUB_REPOSITORY:-}" ]]; then
		echo "${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}"
		return
	fi

	git remote get-url origin | sed -E \
		-e 's#^git@([^:]+):#https://\1/#' \
		-e 's#^ssh://git@([^/]+)/#https://\1/#' \
		-e 's#\.git$##'
}

# ─── next-version ─────────────────────────────────────────────────────────────

next_version() {
	local bump="$1" latest major minor patch

	case "$bump" in
		major | minor | patch) ;;
		*) die "unknown bump '${bump}', expected major, minor or patch" ;;
	esac

	latest="$(latest_tag)"
	if [[ -z "$latest" ]]; then
		echo "1.0.0"
		return
	fi

	IFS=. read -r major minor patch <<<"${latest#v}"
	case "$bump" in
		major) echo "$((major + 1)).0.0" ;;
		minor) echo "${major}.$((minor + 1)).0" ;;
		patch) echo "${major}.${minor}.$((patch + 1))" ;;
	esac
}

# ─── changelog ────────────────────────────────────────────────────────────────

added="" changed="" deprecated="" removed="" fixed="" security=""

# collect_changes sorts the Conventional Commits in <latest tag>..ref into sections.
collect_changes() {
	local ref="$1" prev range subject body type scope bang desc entry
	local re='^([a-zA-Z]+)(\(([^)]*)\))?(!)?:[[:space:]]+(.+)$'

	prev="$(latest_tag)"
	range="$ref"
	[[ -n "$prev" ]] && range="${prev}..${ref}"

	while IFS=$'\x1f' read -r -d $'\x1e' subject body; do
		subject="${subject#$'\n'}"
		[[ "$subject" =~ $re ]] || continue

		type="$(printf '%s' "${BASH_REMATCH[1]}" | tr '[:upper:]' '[:lower:]')"
		scope="${BASH_REMATCH[3]}"
		bang="${BASH_REMATCH[4]}"
		desc="${BASH_REMATCH[5]}"

		# Capitalize the first letter and end with a period, as Keep a Changelog entries do.
		desc="$(printf '%s' "${desc:0:1}" | tr '[:lower:]' '[:upper:]')${desc:1}"
		[[ "$desc" =~ [.!?]$ ]] || desc="${desc}."

		entry="- "
		if [[ -n "$bang" ]] || grep -qE '^BREAKING[ -]CHANGE:' <<<"$body"; then
			entry+="💥 **BREAKING** "
		fi
		[[ -n "$scope" && "$scope" != "security" ]] && entry+="**${scope}:** "
		entry+="${desc}"$'\n'

		if [[ "$scope" == "security" ]]; then
			security+="$entry"
			continue
		fi

		case "$type" in
			feat) added+="$entry" ;;
			refactor | perf | rename) changed+="$entry" ;;
			deprecate) deprecated+="$entry" ;;
			remove) removed+="$entry" ;;
			fix) fixed+="$entry" ;;
			security) security+="$entry" ;;
		esac
	done < <(git log --reverse --no-merges --format='%s%x1f%b%x1e' "$range")
}

print_section() {
	[[ -n "$2" ]] && printf '\n### %s\n\n%s' "$1" "$2"
	return 0
}

render_changelog() {
	printf '## [%s] - %s\n' "$1" "$(date -u +%F)"
	print_section "✨ Added" "$added"
	print_section "♻️ Changed" "$changed"
	print_section "⚠️ Deprecated" "$deprecated"
	print_section "🗑️ Removed" "$removed"
	print_section "🐛 Fixed" "$fixed"
	print_section "🔒 Security" "$security"

	if [[ -z "${added}${changed}${deprecated}${removed}${fixed}${security}" ]]; then
		print_section "♻️ Changed" "- 🧹 Maintenance release with no user-facing changes."$'\n'
	fi
}

# write_changelog inserts the rendered section below "## [Unreleased]" in CHANGELOG.md
# and rewrites the [Unreleased] compare link plus a link for the new version.
write_changelog() {
	local version="$1" prev url version_link unreleased_link section_file out_file

	[[ -f "$CHANGELOG" ]] || die "${CHANGELOG} not found"
	grep -q '^## \[Unreleased\]' "$CHANGELOG" || die "${CHANGELOG} has no '## [Unreleased]' heading"
	! grep -qF "## [${version}]" "$CHANGELOG" || die "${CHANGELOG} already has a ${version} section"

	prev="$(latest_tag)"
	url="$(repo_url)"
	if [[ -n "$prev" ]]; then
		version_link="[${version}]: ${url}/compare/${prev}...v${version}"
	else
		version_link="[${version}]: ${url}/releases/tag/v${version}"
	fi
	unreleased_link="[Unreleased]: ${url}/compare/v${version}...HEAD"

	section_file="$(mktemp)"
	out_file="$(mktemp)"
	render_changelog "$version" >"$section_file"

	awk -v sec="$section_file" -v unreleased="$unreleased_link" -v ver="$version_link" '
		/^## \[Unreleased\]/ {
			print
			print ""
			while ((getline line < sec) > 0) print line
			next
		}
		/^\[Unreleased\]: / { next }
		/^\[[0-9]+\.[0-9]+\.[0-9]+\]: / && !linked {
			print unreleased
			print ver
			linked = 1
		}
		{ print }
		END {
			if (!linked) {
				print ""
				print unreleased
				print ver
			}
		}
	' "$CHANGELOG" >"$out_file"

	cat "$out_file" >"$CHANGELOG"
	rm -f "$section_file" "$out_file"
}

cmd_changelog() {
	local version="" ref="HEAD" write=false

	while (($#)); do
		case "$1" in
			--ref)
				ref="${2:?--ref needs a value}"
				shift 2
				;;
			--write)
				write=true
				shift
				;;
			-*) usage ;;
			*)
				version="$1"
				shift
				;;
		esac
	done
	[[ -n "$version" ]] || usage

	collect_changes "$ref"
	render_changelog "$version"
	if $write; then
		write_changelog "$version"
	fi
}

# ─── build ────────────────────────────────────────────────────────────────────

cmd_build() {
	local version="${1:-}" out="${2:-dist}" target goos goarch name staging binary
	[[ -n "$version" ]] || usage

	rm -rf "$out"
	mkdir -p "$out"

	for target in $BUILD_TARGETS; do
		goos="${target%/*}"
		goarch="${target#*/}"
		name="linkwarden_${version}_${goos}_${goarch}"
		staging="${out}/${name}"
		binary="linkwarden"
		[[ "$goos" == "windows" ]] && binary="linkwarden.exe"

		echo "🔨 ${goos}/${goarch}"
		mkdir -p "$staging"
		CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
			go build -trimpath -ldflags "-s -w -X ${VERSION_PKG}=${version}" -o "${staging}/${binary}" .
		cp README.md "$staging/"
		[[ -f "$CHANGELOG" ]] && cp "$CHANGELOG" "$staging/"

		if [[ "$goos" == "windows" ]]; then
			(cd "$out" && zip -qr "${name}.zip" "$name")
		else
			tar -C "$out" -czf "${out}/${name}.tar.gz" "$name"
		fi
		rm -rf "$staging"
	done

	if command -v sha256sum >/dev/null; then
		(cd "$out" && sha256sum -- *.tar.gz *.zip >checksums.txt)
	else
		(cd "$out" && shasum -a 256 -- *.tar.gz *.zip >checksums.txt)
	fi

	echo "📦 Artifacts in ${out}/:"
	ls -1 "$out"
}

# ─── trigger (default) ────────────────────────────────────────────────────────

cmd_trigger() {
	local bump="" assume_yes=false arg current next unpushed preview answer previous_run run_id

	for arg in "$@"; do
		case "$arg" in
			major | minor | patch) bump="$arg" ;;
			-y | --yes) assume_yes=true ;;
			*) usage ;;
		esac
	done
	[[ -n "$bump" ]] || usage

	command -v gh >/dev/null || die "GitHub CLI not found. Install it from https://cli.github.com"
	gh auth status >/dev/null 2>&1 || die "GitHub CLI isn't authenticated. Run: gh auth login"

	echo "🔄 Fetching origin/main and tags…"
	git fetch --quiet --tags origin main

	current="$(latest_tag)"
	next="$(next_version "$bump")"

	unpushed="$(git rev-list --count origin/main..HEAD 2>/dev/null || echo 0)"
	if [[ "$unpushed" -gt 0 ]]; then
		echo "⚠️  ${unpushed} local commit(s) aren't on origin/main and won't be in this release."
	fi

	preview="$(collect_changes origin/main && render_changelog "$next")"

	echo
	echo "🏷️  ${current:-no releases yet} → v${next} (${bump})"
	echo
	echo "📝 Changelog preview:"
	echo
	# shellcheck disable=SC2001 # indent every line, not a single substitution
	echo "$preview" | sed 's/^/   /'
	echo

	if [[ "$bump" != "major" ]] && grep -q 'BREAKING' <<<"$preview"; then
		echo "💥 This release contains breaking changes. SemVer calls for a major bump."
		echo
	fi

	if ! $assume_yes; then
		read -r -p "🚀 Release v${next}? [y/N] " answer
		[[ "$answer" =~ ^[Yy]$ ]] || {
			echo "🛑 Aborted."
			exit 1
		}
	fi

	previous_run="$(gh run list --workflow "$WORKFLOW" --limit 1 --json databaseId -q '.[0].databaseId' 2>/dev/null || true)"

	gh workflow run "$WORKFLOW" --ref main -f bump="$bump"
	echo "✅ Release workflow triggered."

	# The new run takes a moment to appear; wait for an ID different from the previous one.
	run_id=""
	for _ in 1 2 3 4 5 6 7 8 9 10; do
		sleep 2
		run_id="$(gh run list --workflow "$WORKFLOW" --limit 1 --json databaseId -q '.[0].databaseId' 2>/dev/null || true)"
		[[ -n "$run_id" && "$run_id" != "$previous_run" ]] && break
		run_id=""
	done

	if [[ -n "$run_id" ]]; then
		echo "🔗 $(gh run view "$run_id" --json url -q .url)"
		echo "👀 Follow it with: gh run watch ${run_id}"
	else
		echo "🔗 $(repo_url)/actions/workflows/${WORKFLOW}"
	fi
}

# ─── main ─────────────────────────────────────────────────────────────────────

case "${1:-}" in
	next-version)
		[[ $# -eq 2 ]] || usage
		next_version "$2"
		;;
	changelog)
		shift
		cmd_changelog "$@"
		;;
	build)
		shift
		cmd_build "$@"
		;;
	"" | -h | --help) usage ;;
	*) cmd_trigger "$@" ;;
esac
