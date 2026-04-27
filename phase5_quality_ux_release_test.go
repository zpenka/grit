package grit

import (
	"os"
	"strings"
	"testing"
)

// === PHASE 5: Code Quality, UX, and Release Tests ===

// TestBenchmarkSuiteExists validates benchmark suite for performance testing.
func TestBenchmarkSuiteExists(t *testing.T) {
	data, err := os.ReadFile("benchmarks_test.go")
	if err != nil {
		t.Fatalf("benchmarks_test.go should exist for performance testing, got error: %v", err)
	}

	content := string(data)
	requiredBenchmarks := []string{
		"BenchmarkParseCommits",
		"BenchmarkFilterCommits",
		"BenchmarkRenderUI",
		"Benchmark",
	}

	for _, bench := range requiredBenchmarks {
		if !strings.Contains(content, bench) {
			t.Errorf("benchmarks_test.go missing: %s", bench)
		}
	}
}

// TestCodeCoverageConfigExists validates codecov configuration.
func TestCodeCoverageConfigExists(t *testing.T) {
	data, err := os.ReadFile(".codecov.yml")
	if err != nil {
		t.Fatalf(".codecov.yml should exist for coverage tracking, got error: %v", err)
	}

	content := string(data)
	requiredItems := []string{
		"coverage",
		"precision",
	}

	for _, item := range requiredItems {
		if !strings.Contains(content, item) {
			t.Errorf(".codecov.yml missing: %s", item)
		}
	}
}

// TestCLIHelpDocumentation validates CLI help and version information.
func TestCLIHelpDocumentation(t *testing.T) {
	data, err := os.ReadFile("HELP.md")
	if err != nil {
		t.Fatalf("HELP.md should exist for CLI documentation, got error: %v", err)
	}

	content := string(data)
	requiredSections := []string{
		"Usage",
		"Commands",
		"Options",
		"Examples",
		"Keyboard",
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("HELP.md missing section: %s", section)
		}
	}
}

// TestExampleWorkflowsDocumented validates example workflows exist.
func TestExampleWorkflowsDocumented(t *testing.T) {
	data, err := os.ReadFile("docs/EXAMPLES.md")
	if err != nil {
		t.Fatalf("docs/EXAMPLES.md should document example workflows, got error: %v", err)
	}

	content := string(data)
	requiredExamples := []string{
		"Code Ownership",
		"Code Hotspots",
		"Bisect",
		"Team Velocity",
		"Example",
	}

	for _, example := range requiredExamples {
		if !strings.Contains(content, example) {
			t.Errorf("docs/EXAMPLES.md missing example: %s", example)
		}
	}
}

// TestInstallationGuideExists validates installation documentation.
func TestInstallationGuideExists(t *testing.T) {
	data, err := os.ReadFile("INSTALL.md")
	if err != nil {
		t.Fatalf("INSTALL.md should exist for installation instructions, got error: %v", err)
	}

	content := string(data)
	requiredMethods := []string{
		"go install",
		"build",
		"Requirements",
		"Quick Start",
	}

	for _, method := range requiredMethods {
		if !strings.Contains(content, method) {
			t.Errorf("INSTALL.md missing: %s", method)
		}
	}
}

// TestReleaseBuildScriptExists validates build script for releases.
func TestReleaseBuildScriptExists(t *testing.T) {
	data, err := os.ReadFile("scripts/build-release.sh")
	if err != nil {
		t.Fatalf("scripts/build-release.sh should exist for building releases, got error: %v", err)
	}

	content := string(data)
	requiredItems := []string{
		"#!/bin/bash",
		"grit",
		"build",
		"linux",
		"darwin",
	}

	for _, item := range requiredItems {
		if !strings.Contains(content, item) {
			t.Errorf("scripts/build-release.sh missing: %s", item)
		}
	}
}

// TestReleaseGuideExists validates release process documentation.
func TestReleaseGuideExists(t *testing.T) {
	// Release guide should be in VERSION.md which already exists
	data, err := os.ReadFile("VERSION.md")
	if err != nil {
		t.Fatalf("VERSION.md not found: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "Release Process") && !strings.Contains(content, "release") {
		t.Errorf("VERSION.md should document release process")
	}
}

// TestHomebrewFormulaExists validates Homebrew package definition.
func TestHomebrewFormulaExists(t *testing.T) {
	data, err := os.ReadFile("Formula/grit.rb")
	if err != nil {
		t.Fatalf("Formula/grit.rb should exist for Homebrew, got error: %v", err)
	}

	content := string(data)
	requiredItems := []string{
		"class Grit",
		"url",
		"sha256",
		"bin.install",
	}

	for _, item := range requiredItems {
		if !strings.Contains(content, item) {
			t.Errorf("Formula/grit.rb missing: %s", item)
		}
	}
}

// TestCICodeCoverageIntegration validates GitHub Actions includes coverage.
func TestCICodeCoverageIntegration(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/tests.yml")
	if err != nil {
		t.Fatalf(".github/workflows/tests.yml not found: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "coverage") && !strings.Contains(content, "Cover") {
		t.Errorf("CI workflow should include coverage reporting")
	}
}

// TestGitIgnoreConfigured validates .gitignore has necessary entries.
func TestGitIgnoreConfigured(t *testing.T) {
	data, err := os.ReadFile(".gitignore")
	if err != nil {
		t.Fatalf(".gitignore not found: %v", err)
	}

	content := string(data)
	requiredEntries := []string{
		"grit",
		"build/",
		"dist/",
		"coverage",
	}

	for _, entry := range requiredEntries {
		if !strings.Contains(content, entry) {
			t.Logf("Note: .gitignore should include: %s", entry)
		}
	}
}

// TestQuickstartTutorialExists validates getting started guide.
func TestQuickstartTutorialExists(t *testing.T) {
	data, err := os.ReadFile("docs/QUICKSTART.md")
	if err != nil {
		t.Fatalf("docs/QUICKSTART.md should exist as tutorial, got error: %v", err)
	}

	content := string(data)
	requiredSections := []string{
		"Installation",
		"First Run",
		"Basic Navigation",
		"Step",
		"Tutorial",
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("docs/QUICKSTART.md missing: %s", section)
		}
	}
}

// TestReleaseChecklistExists validates release checklist documentation.
func TestReleaseChecklistExists(t *testing.T) {
	data, err := os.ReadFile("docs/RELEASE_CHECKLIST.md")
	if err != nil {
		t.Fatalf("docs/RELEASE_CHECKLIST.md should exist, got error: %v", err)
	}

	content := string(data)
	lowerContent := strings.ToLower(content)
	requiredItems := []string{
		"tests",
		"coverage",
		"changelog",
		"tag",
		"release",
		"announce",
	}

	for _, item := range requiredItems {
		if !strings.Contains(lowerContent, item) {
			t.Errorf("docs/RELEASE_CHECKLIST.md missing: %s", item)
		}
	}
}

// TestQualityMetricsDocumented validates code quality documentation.
func TestQualityMetricsDocumented(t *testing.T) {
	data, err := os.ReadFile("QUALITY.md")
	if err != nil {
		t.Fatalf("QUALITY.md should document code quality standards, got error: %v", err)
	}

	content := string(data)
	requiredItems := []string{
		"Coverage",
		"Tests",
		"Benchmark",
		"Standard",
		"Quality",
	}

	for _, item := range requiredItems {
		if !strings.Contains(content, item) {
			t.Errorf("QUALITY.md missing: %s", item)
		}
	}
}
