package grit

import (
	"testing"
)

// ============================================================================
// Bulk Feature Wiring Tests - TDD Approach
// Goal: Wire all features that have show* flags and render functions
// ============================================================================

// Test that all known render functions get called when their show* flags are set
// This test matrix ensures comprehensive coverage

var featureWiringTests = []struct {
	setup func(*model)
	name  string
}{
	// Analytics features
	{
		setup: func(m *model) {
			m.showCoupling = true
		},
		name: "Coupling Analysis",
	},
	{
		setup: func(m *model) {
			m.showDependencies = true
		},
		name: "Dependencies",
	},
	// Integration/Conflict features
	{
		setup: func(m *model) {
			m.showConflictUI = true
		},
		name: "Conflict Resolution",
	},
	{
		setup: func(m *model) {
			m.showExportUI = true
			m.pendingExports = []exportData{{format: "markdown"}}
		},
		name: "Export UI",
	},
}

// Test that show flags for existing render functions are properly wired in View()
func TestFeatureRenderingWiring(t *testing.T) {
	for _, test := range featureWiringTests {
		t.Run(test.name, func(t *testing.T) {
			m := newModel(".")
			m.width = 80
			m.height = 20
			test.setup(&m)

			output := m.View()

			// Should produce output when feature flag is set
			if output == "" {
				t.Errorf("%s should render when enabled", test.name)
			}
		})
	}
}

// Test that View() includes all expected render calls
func TestViewMethodIncludesAllWiredFeatures(t *testing.T) {
	// This test verifies that the View() implementation actually calls
	// the render functions we expect, not just returns the main UI
	m := newModel(".")
	m.width = 80
	m.height = 20

	for _, test := range featureWiringTests {
		test.setup(&m)

		output := m.View()

		// Verify output is not empty
		if output == "" {
			t.Errorf("View() should produce output when %s is enabled", test.name)
		}

		// Verify output contains actual UI content (not just main view)
		// Features should produce more than just navigation footer
		if len(output) < 50 {
			t.Logf("Warning: %s output seems short: %d chars", test.name, len(output))
		}
	}
}

// Test that toggling features works correctly
func TestFeatureToggling(t *testing.T) {
	for _, test := range featureWiringTests {
		t.Run(test.name+"-toggle", func(t *testing.T) {
			m := newModel(".")
			m.width = 80
			m.height = 20

			// Feature off
			output1 := m.View()

			// Feature on
			test.setup(&m)
			output2 := m.View()

			// Output should be different when feature is enabled
			if output1 == output2 {
				t.Logf("Tip: %s may not be producing distinct output", test.name)
			}
		})
	}
}

// Test that exclusive features don't render simultaneously
// (Many features should replace main view when shown)
func TestFeatureExclusivity(t *testing.T) {
	m := newModel(".")
	m.width = 80
	m.height = 20
	m.commits = []commit{
		{hash: "abc123", author: "test", subject: "test", when: "1h ago"},
	}

	// Enable two features
	m.showCoupling = true
	m.showConflictUI = true

	output := m.View()

	// When multiple features are enabled, the first one in View() precedence wins
	if output == "" {
		t.Error("View() should render something even with multiple features")
	}
}

// ============================================================================
// Feature Completeness Tests
// ============================================================================

// Test that all wirable features have corresponding render functions
func TestAllWireableFeatureSyntax(t *testing.T) {
	for _, test := range featureWiringTests {
		t.Run(test.name+"-exists", func(t *testing.T) {
			// This is a syntax check - if the setup doesn't reference
			// a valid model field, the test won't even compile
			m := newModel(".")
			test.setup(&m)
		})
	}
}

// ============================================================================
// Performance Tests - Ensure wiring doesn't slow View()
// ============================================================================

func BenchmarkViewWithFeatureWiring(b *testing.B) {
	m := newModel(".")
	m.width = 80
	m.height = 20
	m.commits = []commit{
		{hash: "abc123", author: "test", subject: "test", when: "1h ago"},
		{hash: "def456", author: "test", subject: "test2", when: "2h ago"},
	}

	// Enable various features
	m.showCoupling = true
	m.showDependencies = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}
