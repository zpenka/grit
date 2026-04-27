package grit

import (
	"os"
	"strings"
	"testing"
)

// TestVersionIsOnePointOh verifies the version constant is set to v1.0.0 for release.
func TestVersionIsOnePointOh(t *testing.T) {
	if Version != "1.0.0" {
		t.Errorf("Version should be 1.0.0 for stable release, got %q", Version)
	}
}

// TestChangelogHasV1Release verifies the CHANGELOG documents v1.0.0 release.
func TestChangelogHasV1Release(t *testing.T) {
	data, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Fatalf("CHANGELOG.md not found: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "[1.0.0]") {
		t.Errorf("CHANGELOG.md should document [1.0.0] release section")
	}
}

// TestChangelogHasPhase5 verifies the CHANGELOG documents Phase 5 work.
func TestChangelogHasPhase5(t *testing.T) {
	data, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Fatalf("CHANGELOG.md not found: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Phase 5") {
		t.Errorf("CHANGELOG.md should document Phase 5 improvements")
	}
}
