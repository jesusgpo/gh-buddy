package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// CurrentBranch returns the name of the current git branch.
func CurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// CreateAndCheckout creates a new branch from the current HEAD and checks it out.
func CreateAndCheckout(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch %q: %w", branchName, err)
	}
	return nil
}

// PushBranch pushes the given branch to the remote, setting the upstream.
func PushBranch(remote, branch string) error {
	cmd := exec.Command("git", "push", "-u", remote, branch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to push branch %q to %q: %w", branch, remote, err)
	}
	return nil
}

// HasUncommittedChanges returns true if the working tree has uncommitted changes.
func HasUncommittedChanges() (bool, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}
	return len(strings.TrimSpace(string(out))) > 0, nil
}

// DefaultBranch returns the default branch of the repository (main or master).
func DefaultBranch() (string, error) {
	out, err := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD", "--short").Output()
	if err != nil {
		// Fallback: try common names
		for _, name := range []string{"main", "master"} {
			if err2 := exec.Command("git", "rev-parse", "--verify", "origin/"+name).Run(); err2 == nil {
				return name, nil
			}
		}
		return "", fmt.Errorf("failed to determine default branch: %w", err)
	}
	branch := strings.TrimSpace(string(out))
	// Remove "origin/" prefix
	parts := strings.SplitN(branch, "/", 2)
	if len(parts) == 2 {
		return parts[1], nil
	}
	return branch, nil
}

// RepoSlug returns the "owner/repo" slug parsed from the origin remote URL.
func RepoSlug() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get origin remote URL: %w", err)
	}
	url := strings.TrimSpace(string(out))
	return parseRepoSlug(url)
}

func parseRepoSlug(rawURL string) (string, error) {
	rawURL = strings.TrimSuffix(rawURL, ".git")

	// SSH: git@github.com:owner/repo
	if strings.HasPrefix(rawURL, "git@") {
		parts := strings.SplitN(rawURL, ":", 2)
		if len(parts) == 2 {
			return parts[1], nil
		}
	}

	// HTTPS: https://github.com/owner/repo
	for _, prefix := range []string{"https://", "http://"} {
		if strings.HasPrefix(rawURL, prefix) {
			trimmed := strings.TrimPrefix(rawURL, prefix)
			// Remove host
			idx := strings.Index(trimmed, "/")
			if idx >= 0 {
				return trimmed[idx+1:], nil
			}
		}
	}

	return "", fmt.Errorf("unable to parse repo slug from URL: %s", rawURL)
}

// FetchLatest fetches the latest changes from the remote.
func FetchLatest(remote string) error {
	cmd := exec.Command("git", "fetch", remote)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to fetch from %q: %w", remote, err)
	}
	return nil
}


// SetUpstreamTracking configures the local branch to track the remote branch.
func SetUpstreamTracking(remote, branch string) error {
	if err := exec.Command("git", "fetch", remote, branch).Run(); err != nil {
		return fmt.Errorf("failed to fetch %s/%s: %w", remote, branch, err)
	}
	upstream := fmt.Sprintf("%s/%s", remote, branch)
	if err := exec.Command("git", "branch", "--set-upstream-to="+upstream, branch).Run(); err != nil {
		return fmt.Errorf("failed to set upstream tracking to %s: %w", upstream, err)
	}
	return nil
}

// CreateBranchFrom creates a new branch from a given base branch and checks it out.
func CreateBranchFrom(branchName, baseBranch, remote string) error {
	// Fetch the specific base branch to ensure the remote ref is up to date.
	fetchCmd := exec.Command("git", "fetch", remote, baseBranch)
	if out, err := fetchCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to fetch %s/%s: %w\n%s", remote, baseBranch, err, strings.TrimSpace(string(out)))
	}

	ref := fmt.Sprintf("%s/%s", remote, baseBranch)
	cmd := exec.Command("git", "checkout", "-b", branchName, ref)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create branch %q from %q: %w\n%s", branchName, ref, err, strings.TrimSpace(string(out)))
	}
	return nil
}
