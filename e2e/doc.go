// Package e2e runs the built CLI binary against a real Linkwarden instance.
// The tests are behind the e2e build tag, so a plain `go test ./...` skips them.
//
// Start a local instance and point the tests at it:
//
//	scripts/dev-linkwarden.sh up
//	eval "$(scripts/dev-linkwarden.sh env)"
//	go test -tags e2e -count=1 ./e2e/...
//
// Every test creates uniquely named resources and deletes them afterwards, so the
// suite can run against a long-lived dev instance. It refuses to run against
// anything but localhost, so it can never touch a real account.
package e2e
