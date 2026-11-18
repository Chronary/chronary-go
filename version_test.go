package chronary

import (
	"os"
	"strings"
	"testing"
)

// VERSION and Version must be equal: the release coordinator refuses to
// publish otherwise (#2397). The PR gate enforces this through
// apps/api/test/unit/shared/release-version-drift.test.ts; this test keeps the
// rule in the public mirror, whose CI runs `go test ./...`.
func TestVersionMatchesVersionFile(t *testing.T) {
	raw, err := os.ReadFile("VERSION")
	if err != nil {
		t.Fatalf("read VERSION: %v", err)
	}
	if want := strings.TrimSpace(string(raw)); Version != want {
		t.Fatalf("version.go Version = %q, VERSION file = %q", Version, want)
	}
}
