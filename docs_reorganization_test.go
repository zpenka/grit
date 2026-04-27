package grit

import (
	"os"
	"strings"
	"testing"
)

// TestDocumentationIndexExists validates docs/README.md serves as documentation index.
func TestDocumentationIndexExists(t *testing.T) {
	data, err := os.ReadFile("docs/README.md")
	if err != nil {
		t.Fatalf("docs/README.md should exist as documentation index, got error: %v", err)
	}

	content := string(data)
	requiredSections := []string{
		"Core Infrastructure",
		"Rendering",
		"Features",
		"Table of Contents",
		"modules/",
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("docs/README.md missing section: %s", section)
		}
	}
}

// TestDocsNavigationFileExists validates DOCS.md provides overall documentation navigation.
func TestDocsNavigationFileExists(t *testing.T) {
	data, err := os.ReadFile("DOCS.md")
	if err != nil {
		t.Fatalf("DOCS.md should exist as documentation navigation, got error: %v", err)
	}

	content := string(data)
	requiredItems := []string{
		"README",
		"DEVELOPER",
		"CONTRIBUTING",
		"ARCHITECTURE",
		"modules/",
		"Quick Start",
	}

	for _, item := range requiredItems {
		if !strings.Contains(content, item) {
			t.Errorf("DOCS.md missing reference to: %s", item)
		}
	}
}

// TestModuleDocumentationIndexed validates all modules are documented in docs/README.md.
func TestModuleDocumentationIndexed(t *testing.T) {
	modules := []string{
		"PARSING.md",
		"NAVIGATION.md",
		"FILTERING.md",
		"CACHING.md",
		"OPTIMIZATION.md",
		"RENDERING.md",
		"RENDER_CONSOLIDATION.md",
		"ANALYTICS.md",
		"GIT_OPS.md",
		"WORKFLOWS.md",
		"VISUALIZATION.md",
		"INTEGRATION.md",
		"TEAM_AI.md",
	}

	data, err := os.ReadFile("docs/README.md")
	if err != nil {
		t.Fatalf("docs/README.md not found: %v", err)
	}

	content := string(data)
	for _, module := range modules {
		if !strings.Contains(content, module) {
			t.Errorf("docs/README.md missing module: %s", module)
		}
	}
}

// TestProgressAndChangelogConsolidation validates progress tracking information exists.
func TestProgressAndChangelogConsolidation(t *testing.T) {
	// Check that phase information is consolidated properly
	data, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Fatalf("CHANGELOG.md not found: %v", err)
	}

	content := string(data)
	phases := []string{"Phase 1", "Phase 2", "Phase 3"}

	for _, phase := range phases {
		if !strings.Contains(content, phase) {
			t.Errorf("CHANGELOG.md missing %s", phase)
		}
	}
}

// TestNoOrphanedDocumentationFiles checks for old/unused documentation.
func TestNoOrphanedDocumentationFiles(t *testing.T) {
	orphanCandidates := []string{
		"REFACTORING.md",
		"REFACTORING_PLAN.md",
		"PHASE_4C_PLAN.md",
		"ARCHITECTURE_COMPREHENSIVE.md",
	}

	for _, orphan := range orphanCandidates {
		if _, err := os.Stat(orphan); err == nil {
			t.Errorf("Orphaned file should be removed: %s", orphan)
		}
	}
}

// TestDocumentationStructure validates overall documentation hierarchy.
func TestDocumentationStructure(t *testing.T) {
	expectedFiles := map[string]bool{
		"README.md":                     true,
		"DOCS.md":                       true,
		"CLAUDE.md":                     true,
		"ARCHITECTURE.md":               true,
		"DEVELOPER.md":                  true,
		"CONTRIBUTING.md":               true,
		"VERSION.md":                    true,
		"CHANGELOG.md":                  true,
		"docs/README.md":                true,
		".github/pull_request_template.md": true,
		"docs/modules/PARSING.md":       true,
		"docs/modules/ANALYTICS.md":     true,
	}

	for file := range expectedFiles {
		if _, err := os.Stat(file); err != nil {
			t.Errorf("Expected documentation file missing: %s", file)
		}
	}
}

// TestDocumentationLinking validates that docs reference each other properly.
func TestDocumentationLinking(t *testing.T) {
	// README should link to major docs
	readmeData, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("README.md not found: %v", err)
	}

	readmeContent := string(readmeData)
	readmeLinks := []string{
		"DEVELOPER.md",
		"CONTRIBUTING.md",
		"docs/modules/",
	}

	for _, link := range readmeLinks {
		if !strings.Contains(readmeContent, link) {
			t.Errorf("README.md should link to: %s", link)
		}
	}

	// DOCS.md should be comprehensive navigation
	docsData, err := os.ReadFile("DOCS.md")
	if err != nil {
		t.Fatalf("DOCS.md not found: %v", err)
	}

	docsContent := string(docsData)
	docsLinks := []string{
		"README",
		"DEVELOPER",
		"CONTRIBUTING",
		"ARCHITECTURE",
		"VERSION",
		"CHANGELOG",
		"docs/modules",
	}

	for _, link := range docsLinks {
		if !strings.Contains(docsContent, link) {
			t.Errorf("DOCS.md should link to: %s", link)
		}
	}
}
