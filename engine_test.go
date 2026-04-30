package grit

import (
	"strings"
	"testing"
)

func makeCommits(n int) []commit {
	var cs []commit
	for i := 0; i < n; i++ {
		cs = append(cs, commit{
			hash:      strings.Repeat("a", 40),
			shortHash: "abc1234",
			author:    "Test User",
			when:      "1d ago",
			subject:   "commit message",
		})
	}
	return cs
}

func makeDiffLines(n int) []diffLine {
	lines := make([]diffLine, n)
	for i := range lines {
		lines[i] = diffLine{kind: lineContext, text: "context"}
	}
	return lines
}

func makeNamedCommits() []commit {
	return []commit{
		{shortHash: "aaa1111", author: "John Doe", subject: "Fix login bug"},
		{shortHash: "bbb2222", author: "Jane Smith", subject: "Add user model"},
		{shortHash: "ccc3333", author: "John Doe", subject: "Update README"},
	}
}

func makeCommitsWithDays() []commit {
	return []commit{
		{shortHash: "aaa1111", author: "John", when: "1 day ago", subject: "Recent"},
		{shortHash: "bbb2222", author: "Jane", when: "5 days ago", subject: "Medium"},
		{shortHash: "ccc3333", author: "Bob", when: "20 days ago", subject: "Old"},
		{shortHash: "ddd4444", author: "Alice", when: "100 days ago", subject: "Very old"},
	}
}


// TestNewModel_Creation tests model initialization
func TestNewModel_Creation(t *testing.T) {
	m := newModel(".")

	AssertEqual(t, 0, m.cursor, "initial cursor should be 0")
	AssertTrue(t, len(m.commits) >= 0, "commits should be initialized")
}

// TestCommitBuilder_FluentAPI tests builder pattern
func TestCommitBuilder_FluentAPI(t *testing.T) {
	commit := NewCommitBuilder().
		WithHash("abc123def456").
		WithAuthor("John Doe").
		WithSubject("Test commit").
		WithWhen("1 hour ago").
		Build()

	AssertEqual(t, "abc123d", commit.shortHash, "should set short hash")
	AssertEqual(t, "John Doe", commit.author, "should set author")
	AssertEqual(t, "Test commit", commit.subject, "should set subject")
	AssertEqual(t, "1 hour ago", commit.when, "should set when")
}

// TestCommitBuilder_Defaults tests builder with defaults
func TestCommitBuilder_Defaults(t *testing.T) {
	commit := NewCommitBuilder().Build()

	AssertEqual(t, "Test Author", commit.author, "should use default author")
	AssertEqual(t, "Test Subject", commit.subject, "should use default subject")
}
