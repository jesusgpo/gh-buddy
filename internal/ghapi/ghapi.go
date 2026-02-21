package ghapi

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Issue represents a GitHub issue.
type Issue struct {
	Number int     `json:"number"`
	Title  string  `json:"title"`
	Body   string  `json:"body"`
	Labels []Label `json:"labels"`
	State  string  `json:"state"`
	URL    string  `json:"html_url"`
}

// Label represents a GitHub issue label.
type Label struct {
	Name string `json:"name"`
}

// PullRequest represents a created pull request.
type PullRequest struct {
	Number int    `json:"number"`
	URL    string `json:"html_url"`
	Title  string `json:"title"`
}

// GetIssue fetches details of a GitHub issue by number.
func GetIssue(repo string, number int) (*Issue, error) {
	out, err := exec.Command("gh", "api",
		fmt.Sprintf("repos/%s/issues/%d", repo, number),
	).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue #%d: %w", number, err)
	}
	var issue Issue
	if err := json.Unmarshal(out, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse issue: %w", err)
	}
	return &issue, nil
}

// ListOpenIssues lists open issues assigned to the current user.
func ListOpenIssues(repo string) ([]Issue, error) {
	out, err := exec.Command("gh", "issue", "list",
		"--repo", repo,
		"--assignee", "@me",
		"--state", "open",
		"--json", "number,title,labels,state,url",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}
	var issues []Issue
	if err := json.Unmarshal(out, &issues); err != nil {
		return nil, fmt.Errorf("failed to parse issues: %w", err)
	}
	return issues, nil
}

// CreatePR creates a pull request via the gh CLI.
func CreatePR(repo, title, body, base, head string, draft bool, labels []string) (*PullRequest, error) {
	args := []string{"pr", "create",
		"--repo", repo,
		"--title", title,
		"--body", body,
		"--base", base,
		"--head", head,
	}
	if draft {
		args = append(args, "--draft")
	}
	for _, l := range labels {
		args = append(args, "--label", l)
	}
	out, err := exec.Command("gh", args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %s: %w", string(out), err)
	}
	// gh pr create outputs the PR URL on success
	url := strings.TrimSpace(string(out))
	pr := &PullRequest{URL: url}

	// Try to extract PR number from URL
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		if num, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
			pr.Number = num
		}
	}
	pr.Title = title
	return pr, nil
}

// ListLabels lists available labels for a repository.
func ListLabels(repo string) ([]Label, error) {
	out, err := exec.Command("gh", "api",
		fmt.Sprintf("repos/%s/labels", repo),
		"--paginate",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list labels: %w", err)
	}
	var labels []Label
	if err := json.Unmarshal(out, &labels); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %w", err)
	}
	return labels, nil
}

// CurrentUser returns the currently authenticated GitHub username.
func CurrentUser() (string, error) {
	out, err := exec.Command("gh", "api", "user", "--jq", ".login").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
