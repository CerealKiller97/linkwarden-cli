#!/usr/bin/env bash
# 🚀 Release helper for linkwarden-cli. Compatible with bash 3.2 (macOS's /bin/bash).
#
# Nobody can push to main directly, so a release takes two steps:
#
#   scripts/release.sh prepare <major|minor|patch> [-y|--yes]
#       1. Opens a pull request that adds the next version's section to CHANGELOG.md. The
#          section is built from origin/main; edit it on the PR branch if you like, then merge.
#
#   scripts/release.sh <major|minor|patch> [-y|--yes]
#       2. Triggers the "release" GitHub Actions workflow for origin/main, which tags main's
#          HEAD and publishes. Fails if origin/main's CHANGELOG.md has no section for the
#          next version yet (run step 1 and merge its PR first).
#
#   Both show the next version and a changelog preview, then ask to confirm (-y skips it).
#   Both require the GitHub CLI (https://cli.github.com), authenticated with `gh auth login`.
#
#   The workflow drives the release through these subcommands. Run them locally to preview:
#
#   scripts/release.sh next-version <major|minor|patch>
#       Prints the next version (without "v") from the highest vX.Y.Z tag.
#       With no tags yet, the first release is 1.0.0.
#
#   scripts/release.sh changelog <version> [--ref REF] [--write | --existing]
#       Prints a Keep a Changelog (https://keepachangelog.com/en/1.1.0/) section built from
#       the Conventional Commits in <latest tag>..REF (REF defaults to HEAD). --write also
#       inserts it into CHANGELOG.md below "## [Unreleased]" and updates the compare links.
#       If CHANGELOG.md already has a <version> section (written or edited by hand), --write
#       keeps the file as it is and prints that section instead.
#       --existing prints the <version> section already in CHANGELOG.md, and fails if there
#       is none. The release workflow uses it, since it can't commit to main.
#
#   Binaries and packages are built by GoReleaser (.goreleaser.yaml). Preview them with:
#       goreleaser release --snapshot --clean
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
# marks the entry with 💥 **BREAKING**, except in the first release (no vX.Y.Z tag yet).
set -euo pipefail

cd "$(dirname "$0")/.."

CHANGELOG="CHANGELOG.md"
WORKFLOW="release.yml"
BASE="main"

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
		# Nothing can break before the first release, so only mark entries once a tag exists.
		if [[ -n "$prev" ]] && { [[ -n "$bang" ]] || grep -qE '^BREAKING[ -]CHANGE:' <<<"$body"; }; then
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

# has_section reports whether CHANGELOG.md (or the file given as $2) already has a
# "## [<version>]" section.
has_section() {
	local file="${2:-$CHANGELOG}"
	[[ -f "$file" ]] && grep -qF "## [${1}]" "$file"
}

# existing_section prints the "## [<version>]" section of CHANGELOG.md (or the file given as
# $2), from its heading up to the next "## [" heading or the link definitions, without
# trailing blank lines.
existing_section() {
	awk -v heading="## [${1}]" '
		index($0, heading) == 1 { found = 1; print; next }
		found && (/^## \[/ || /^\[[^]]+\]: /) { exit }
		found { lines[++n] = $0 }
		END {
			while (n > 0 && lines[n] ~ /^[[:space:]]*$/) n--
			for (i = 1; i <= n; i++) print lines[i]
		}
	' "${2:-$CHANGELOG}"
}

# remote_changelog writes origin/main's CHANGELOG.md to a temporary file and prints its path.
remote_changelog() {
	local file
	file="$(mktemp)"
	git show "origin/${BASE}:${CHANGELOG}" >"$file" 2>/dev/null || true
	echo "$file"
}

# require_gh fails unless the GitHub CLI is installed and authenticated.
require_gh() {
	command -v gh >/dev/null || die "GitHub CLI not found. Install it from https://cli.github.com"
	gh auth status >/dev/null 2>&1 || die "GitHub CLI isn't authenticated. Run: gh auth login"
}

# parse_bump sets bump and assume_yes from "<major|minor|patch> [-y|--yes]".
parse_bump() {
	local arg
	bump="" assume_yes=false
	for arg in "$@"; do
		case "$arg" in
			major | minor | patch) bump="$arg" ;;
			-y | --yes) assume_yes=true ;;
			*) usage ;;
		esac
	done
	[[ -n "$bump" ]] || usage
}

# confirm asks "$1 [y/N]" unless assume_yes is set, and exits on anything but y.
confirm() {
	local answer
	$assume_yes && return
	read -r -p "$1 [y/N] " answer
	[[ "$answer" =~ ^[Yy]$ ]] || {
		echo "🛑 Aborted."
		exit 1
	}
}

# show_preview prints the version bump and an indented changelog section.
show_preview() {
	local current="$1" next="$2" preview="$3"
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
	local version="" ref="HEAD" write=false existing=false

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
			--existing)
				existing=true
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
	if $write && $existing; then
		usage
	fi

	if $existing; then
		has_section "$version" ||
			die "${CHANGELOG} has no ${version} section. Run 'scripts/release.sh prepare <bump>' and merge its pull request first."
		existing_section "$version"
		return
	fi

	if $write && has_section "$version"; then
		echo "📝 ${CHANGELOG} already has a ${version} section; using it as is." >&2
		existing_section "$version"
		return
	fi

	collect_changes "$ref"
	render_changelog "$version"
	if $write; then
		write_changelog "$version"
	fi
}

# ─── prepare ──────────────────────────────────────────────────────────────────

cmd_prepare() {
	local bump assume_yes current next remote branch title worktree body

	parse_bump "$@"
	require_gh

	echo "🔄 Fetching origin/${BASE} and tags…"
	git fetch --quiet --tags origin "$BASE"

	current="$(latest_tag)"
	next="$(next_version "$bump")"
	branch="release/v${next}"
	title="chore(release): v${next}"

	remote="$(remote_changelog)"
	if has_section "$next" "$remote"; then
		rm -f "$remote"
		echo "✅ origin/${BASE} already has a ${next} section in ${CHANGELOG}."
		echo "🚀 Release it with: scripts/release.sh ${bump}"
		return
	fi
	rm -f "$remote"

	if git ls-remote --exit-code --heads origin "$branch" >/dev/null 2>&1; then
		die "Branch ${branch} already exists on origin. Merge or close its pull request first."
	fi

	show_preview "$current" "$next" "$(collect_changes "origin/${BASE}" && render_changelog "$next")"
	confirm "📝 Open a pull request adding the v${next} changelog?"

	# Work in a throwaway worktree so the current checkout and its changes stay untouched.
	worktree="$(mktemp -d)"
	git worktree add --quiet -b "$branch" "$worktree" "origin/${BASE}"
	# shellcheck disable=SC2064 # expand $worktree now, while it's set
	trap "git worktree remove --force '${worktree}' 2>/dev/null || true" EXIT

	(
		cd "$worktree"
		collect_changes "origin/${BASE}"
		write_changelog "$next"
		git add "$CHANGELOG"
		git commit --quiet -m "$title"
		git push --quiet -u origin "$branch"
	)

	body="Adds the v${next} section to \`${CHANGELOG}\`, built from the Conventional Commits since ${current:-the first commit}.

Edit the section on this branch if it needs it. Once this is merged, release with:

\`\`\`bash
scripts/release.sh ${bump}
\`\`\`"
	gh pr create --base "$BASE" --head "$branch" --title "$title" --body "$body"
}

# ─── trigger (default) ────────────────────────────────────────────────────────

cmd_trigger() {
	local bump assume_yes current next unpushed remote preview previous_run run_id

	parse_bump "$@"
	require_gh

	echo "🔄 Fetching origin/${BASE} and tags…"
	git fetch --quiet --tags origin "$BASE"

	current="$(latest_tag)"
	next="$(next_version "$bump")"

	# The workflow can't commit to main, so the changelog has to be merged already.
	remote="$(remote_changelog)"
	if ! has_section "$next" "$remote"; then
		rm -f "$remote"
		die "origin/${BASE}'s ${CHANGELOG} has no ${next} section. Run 'scripts/release.sh prepare ${bump}' and merge its pull request first."
	fi
	preview="$(existing_section "$next" "$remote")"
	rm -f "$remote"

	unpushed="$(git rev-list --count "origin/${BASE}..HEAD" 2>/dev/null || echo 0)"
	if [[ "$unpushed" -gt 0 ]]; then
		echo "⚠️  ${unpushed} local commit(s) aren't on origin/${BASE} and won't be in this release."
	fi

	show_preview "$current" "$next" "$preview"
	confirm "🚀 Release v${next}?"

	previous_run="$(gh run list --workflow "$WORKFLOW" --limit 1 --json databaseId -q '.[0].databaseId' 2>/dev/null || true)"

	gh workflow run "$WORKFLOW" --ref "$BASE" -f bump="$bump"
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
	prepare)
		shift
		cmd_prepare "$@"
		;;
	"" | -h | --help) usage ;;
	*) cmd_trigger "$@" ;;
esac
