#!/usr/bin/env bash
# End-to-end tryout of --json and -q against your account:
#   1. creates a test collection and captures its ID with -q,
#   2. imports two links into it, one of them invalid on purpose,
#   3. reports the failed one from the --json summary,
#   4. lists what was saved,
#   5. asks, then deletes the test collection and its links.
#
# If it stops early (an error, or Ctrl+C), it tells you how to delete the
# collection it made.
#
# Usage:
#   examples/try-it.sh
set -euo pipefail
here="$(cd "$(dirname "$0")" && pwd)"
. "$here/lib.sh"

step() {
  printf '\n\033[1m%s\033[0m\n' "$*"
}

col=""
deleted=false
links_file=$(mktemp)

cleanup() {
  rm -f "$links_file"
  if [ -n "$col" ] && [ "$deleted" = false ]; then
    echo
    echo "The test collection $col is still there. Delete it with:"
    echo "  $LW collections delete --id $col"
  fi
}
trap cleanup EXIT

# A unique name, since creating a collection whose name is taken fails.
name="CLI test $(date +%Y-%m-%d\ %H.%M.%S)"

step "1. Create a collection, keeping its ID with -q"
col=$("$LW" collections create --name "$name" -q)
echo "Created \"$name\" with ID $col"

step "2–3. Import two links, one invalid, and report the failure"
cat > "$links_file" <<'EOF'
[
  { "url": "https://go.dev/doc/effective_go", "title": "Effective Go" },
  { "url": "not a url", "title": "Invalid on purpose" }
]
EOF
if "$here/import-and-report.sh" "$links_file" "$col"; then
  echo "(No failures this time.)"
else
  echo "import-and-report.sh exited $?, as a cron job would see it."
fi

step "4. Links now in the collection (--json, shaped with jq)"
"$LW" links list --collection-id "$col" --json | jq -r '.[] | "  \(.id)  \(.name)  \(.url)"'

step "5. Clean up"
read -r -p "Delete collection $col and its links? This can't be undone. [y/N] " answer
case "$answer" in
  y | Y)
    "$LW" collections delete --id "$col" --yes -q
    deleted=true
    echo "Deleted."
    ;;
  *)
    deleted=true # kept on purpose: no reminder needed
    echo "Kept it. Delete it later with: $LW collections delete --id $col"
    ;;
esac
