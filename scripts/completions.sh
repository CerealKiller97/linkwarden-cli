#!/usr/bin/env bash
# 🐚 Generates bash, zsh and fish completion scripts into completions/ for GoReleaser to package.
set -euo pipefail

cd "$(dirname "$0")/.."

rm -rf completions
mkdir -p completions

# The root command insists on LINKWARDEN_TOKEN even for `completion`; any value will do here.
for shell in bash zsh fish; do
	LINKWARDEN_TOKEN=completions go run . completion "$shell" >"completions/linkwarden.${shell}"
done
