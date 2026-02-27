package branch

import (
	"fmt"
	"regexp"
	"strings"
)

// IssueType represents the type of issue for branch naming.
type IssueType string

const (
	Feature  IssueType = "feature"
	Bugfix   IssueType = "bugfix"
	Hotfix   IssueType = "hotfix"
	Release  IssueType = "release"
	Chore    IssueType = "chore"
	Docs     IssueType = "docs"
	Refactor IssueType = "refactor"
	Test     IssueType = "test"
	Internal IssueType = "internal"
)

// AllIssueTypes returns all valid issue types.
func AllIssueTypes() []IssueType {
	return []IssueType{Feature, Bugfix, Hotfix, Release, Chore, Docs, Refactor, Test, Internal}
}

// AllIssueTypeStrings returns all valid issue types as strings.
func AllIssueTypeStrings() []string {
	types := AllIssueTypes()
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = string(t)
	}
	return result
}

// ValidIssueType checks if the given string is a valid issue type.
func ValidIssueType(s string) bool {
	for _, t := range AllIssueTypes() {
		if string(t) == s {
			return true
		}
	}
	return false
}

var nonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// GenerateName generates a branch name from issue type, number, and title.
// Format: <type>/GH-<issue-number>-<slugified-title>
func GenerateName(issueType IssueType, issueNumber int, title string) string {
	slug := slugify(title)
	if issueNumber > 0 {
		return fmt.Sprintf("%s/GH-%d-%s", issueType, issueNumber, slug)
	}
	return fmt.Sprintf("%s/%s", issueType, slug)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	// Limit length
	if len(s) > 60 {
		s = s[:60]
		// Don't end on a hyphen
		s = strings.TrimRight(s, "-")
	}
	return s
}
