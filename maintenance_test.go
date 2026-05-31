package grit

import (
	"os"
	"strings"
	"testing"
)

// TestMaintenancePhaseDocumentation validates Phase 3 maintenance documentation exists and is properly structured.
func TestMaintenancePhaseDocumentation(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		requiredText []string
		shouldExist  bool
	}{
		{
			name:     "CONTRIBUTING.md exists",
			filePath: "CONTRIBUTING.md",
			requiredText: []string{
				"Contributing", "code",
				"Tests", "test",
				"Guidelines", "style",
			},
			shouldExist: true,
		},
		{
			name:     "VERSION.md exists",
			filePath: "VERSION.md",
			requiredText: []string{
				"Version", "Semantic",
				"Release", "tag",
			},
			shouldExist: true,
		},
		{
			name:     "Issue template exists",
			filePath: ".github/ISSUE_TEMPLATE/bug_report.md",
			requiredText: []string{
				"bug", "Bug", "describe",
			},
			shouldExist: true,
		},
		{
			name:     "Feature request template exists",
			filePath: ".github/ISSUE_TEMPLATE/feature_request.md",
			requiredText: []string{
				"feature", "Feature", "request",
			},
			shouldExist: true,
		},
		{
			name:     "PR template exists",
			filePath: ".github/pull_request_template.md",
			requiredText: []string{
				"Summary", "Changes", "Test",
			},
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.filePath)
			if tt.shouldExist && err != nil {
				t.Fatalf("Expected file %s to exist, but got error: %v", tt.filePath, err)
			}
			if !tt.shouldExist && err == nil {
				t.Fatalf("Expected file %s to not exist, but it does", tt.filePath)
			}

			if tt.shouldExist && len(data) > 0 {
				content := string(data)
				lowerContent := strings.ToLower(content)
				for _, required := range tt.requiredText {
					if !strings.Contains(lowerContent, strings.ToLower(required)) {
						t.Errorf("Expected %s to contain '%s' (case-insensitive), but it doesn't", tt.filePath, required)
					}
				}
			}
		})
	}
}

// TestBranchCleanup validates that development branches are documented.
func TestBranchCleanup(t *testing.T) {
	// This test just validates the concept - branches are ephemeral
	// so we just check that cleanup guidance exists in CONTRIBUTING.md
	data, err := os.ReadFile("CONTRIBUTING.md")
	if err != nil {
		t.Fatalf("CONTRIBUTING.md not found: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "branch") && !strings.Contains(content, "cleanup") {
		t.Errorf("CONTRIBUTING.md should document branch management")
	}
}

// TestVersionManagement validates versioning documentation exists.
func TestVersionManagement(t *testing.T) {
	data, err := os.ReadFile("VERSION.md")
	if err != nil {
		t.Fatalf("VERSION.md not found: %v", err)
	}

	content := string(data)
	requiredSections := []string{
		"Semantic Versioning",
		"Version Format",
		"Release",
		"Git Tags",
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("VERSION.md missing section: %s", section)
		}
	}
}
