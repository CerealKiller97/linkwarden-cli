#!/usr/bin/env bash
#
# Fills an empty Linkwarden account with the collections, links and tags the
# demo GIFs show. Run it once against a throwaway instance, with LINKWARDEN_URL
# and LINKWARDEN_TOKEN pointing at it:
#
#   .github/assets/demo/seed.sh
#
# The tapes expect the account to be empty before seeding, so that "Go" gets
# collection ID 1 and "Go Concurrency Patterns: Pipelines" link ID 2.

set -euo pipefail

LINKWARDEN="${LINKWARDEN_BIN:-linkwarden}"

# collection NAME DESCRIPTION COLOR — prints the new collection's ID.
collection() {
	"$LINKWARDEN" collections create --name "$1" --description "$2" --color "$3" -q
}

# link COLLECTION_ID URL NAME [TAG...]
link() {
	local collection_id="$1" url="$2" name="$3"
	shift 3
	local args=()
	for tag in "$@"; do
		args+=(--tag "$tag")
	done
	"$LINKWARDEN" links create "$url" --name "$name" --collection-id "$collection_id" "${args[@]}" -q >/dev/null
}

go_id="$(collection "Go" "Docs, talks and blog posts about Go" "#00add8")"
reading_id="$(collection "Reading List" "Long-form articles for the weekend" "#f59e0b")"
infra_id="$(collection "Infra" "Kubernetes, containers and ops" "#6366f1")"
design_id="$(collection "Design" "Software design and architecture" "#ec4899")"

link "$go_id" https://go.dev/doc/effective_go "Effective Go" golang docs
link "$go_id" https://go.dev/blog/pipelines "Go Concurrency Patterns: Pipelines" golang concurrency
link "$go_id" https://go.dev/ref/spec "The Go Programming Language Specification" golang docs
link "$go_id" https://gobyexample.com "Go by Example" golang tutorial
link "$go_id" https://go.dev/blog/errors-are-values "Errors are values" golang
link "$infra_id" https://k8spatterns.io "Kubernetes Patterns" k8s books
link "$infra_id" https://kubernetes.io/docs/concepts/overview/ "Kubernetes Overview" k8s docs
link "$infra_id" https://docs.docker.com/build/building/multi-stage/ "Multi-stage builds" docker docs
link "$infra_id" https://sre.google/sre-book/table-of-contents/ "Site Reliability Engineering" sre books
link "$design_id" https://12factor.net "The Twelve-Factor App" architecture
link "$design_id" https://martinfowler.com/articles/microservices.html "Microservices" architecture
link "$design_id" https://www.hillelwayne.com/post/what-we-know-we-dont-know/ "What We Know We Don't Know" engineering
link "$reading_id" https://paulgraham.com/greatwork.html "How to Do Great Work" essays
link "$reading_id" https://danluu.com/everything-is-broken/ "Everything is broken" essays
link "$reading_id" https://jvns.ca/blog/2024/11/18/how-to-import-a-javascript-library/ "Importing a frontend JS library" javascript
link "$reading_id" https://www.charmbracelet.com "Charm" terminal tools

echo "Seeded 4 collections and 16 links"
