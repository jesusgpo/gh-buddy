package cmd

import (
	"fmt"
	"strconv"

	"github.com/jesusgpo/gh-buddy/internal/branch"
	"github.com/jesusgpo/gh-buddy/internal/ghapi"
	"github.com/jesusgpo/gh-buddy/internal/git"
	"github.com/jesusgpo/gh-buddy/internal/prompt"
	"github.com/jesusgpo/gh-buddy/internal/ui"
	"github.com/spf13/cobra"
)

func newCreateBranchCmd() *cobra.Command {
	var (
		issueNumber int
		issueType   string
		baseBranch  string
	)

	cmd := &cobra.Command{
		Use:   "create-branch",
		Short: "Create a local branch from an issue",
		Long: `Create a local branch following naming conventions.

If an issue number is provided, the branch name will be generated from the issue title.
The branch type can be one of: feature, bugfix, hotfix, release, chore, docs, refactor, test, internal.`,
		Example: `  # Create a branch from issue #42
  gh buddy create-branch --issue 42

  # Create a branch with a specific type
  gh buddy create-branch --issue 42 --type bugfix

  # Create a branch from a different base
  gh buddy create-branch --issue 42 --base develop

  # Use defaults without prompts
  gh buddy create-branch --issue 42 -y`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateBranch(issueNumber, issueType, baseBranch)
		},
	}

	cmd.Flags().IntVarP(&issueNumber, "issue", "i", 0, "issue number to create the branch from")
	cmd.Flags().StringVarP(&issueType, "type", "t", "", "branch type (feature, bugfix, hotfix, release, chore, docs, refactor, test, internal)")
	cmd.Flags().StringVarP(&baseBranch, "base", "b", "", "base branch to create from (default: repo default branch)")

	return cmd
}

func runCreateBranch(issueNumber int, issueType, baseBranch string) error {
	repo, err := git.RepoSlug()
	if err != nil {
		return fmt.Errorf("not in a git repository or no origin remote: %w", err)
	}

	// If no issue number provided, prompt for selection or manual input
	if issueNumber == 0 {
		issueNumber, err = promptForIssue(repo)
		if err != nil {
			return err
		}
	}

	// Fetch issue details
	issue, err := ghapi.GetIssue(repo, issueNumber)
	if err != nil {
		return err
	}

	ui.IssuePanel(issue.Number, issue.Title)

	// Determine issue type
	if issueType == "" {
		issueType = inferIssueType(issue)
	}

	if !useDefaults {
		types := branch.AllIssueTypeStrings()
		defaultIdx := 0
		for i, t := range types {
			if t == issueType {
				defaultIdx = i
				break
			}
		}
		idx, err := prompt.SelectWithDefault("Select branch type:", types, defaultIdx)
		if err != nil {
			return err
		}
		issueType = types[idx]
	} else if issueType == "" {
		issueType = "feature"
	}

	if !branch.ValidIssueType(issueType) {
		return fmt.Errorf("invalid branch type %q. Valid types: %v", issueType, branch.AllIssueTypeStrings())
	}

	// Determine base branch
	if baseBranch == "" {
		defaultBase, err := git.DefaultBranch()
		if err != nil {
			defaultBase = "main"
		}
		if !useDefaults {
			baseBranch = prompt.Input("Base branch", defaultBase)
		} else {
			baseBranch = defaultBase
		}
	}

	// Generate branch name
	branchName := branch.GenerateName(branch.IssueType(issueType), issueNumber, issue.Title)

	if !useDefaults {
		branchName = prompt.Input("Branch name", branchName)
	}

	ui.BranchPanel(branchName, baseBranch)

	// Create the branch
	if err := git.CreateBranchFrom(branchName, baseBranch, "origin"); err != nil {
		return err
	}

	ui.Success("Branch %q created and checked out successfully!", branchName)

	// Ask to push
	shouldPush := useDefaults || prompt.Confirm("Push branch to origin?", true)
	if shouldPush {
		// Use `gh issue develop` to push the branch to GitHub and link it to the
		// issue in one step. If that fails, fall back to a regular git push.
		if linkErr := ghapi.LinkBranchToIssue(repo, issueNumber, branchName, baseBranch); linkErr != nil {
			ui.Warning("Could not create linked branch via gh issue develop (%v), falling back to git push", linkErr)
			if err := git.PushBranch("origin", branchName); err != nil {
				return err
			}
			ui.Success("Branch pushed to origin")
		} else {
			// Branch now exists on remote; configure local tracking
			if err := git.SetUpstreamTracking("origin", branchName); err != nil {
				ui.Warning("Branch pushed but could not set upstream tracking: %v", err)
			}
			ui.Success("Branch pushed to origin")
			ui.Success("Branch linked to issue #%d", issueNumber)
		}
	}

	return nil
}

func promptForIssue(repo string) (int, error) {
	// List open issues assigned to the user
	issues, err := ghapi.ListOpenIssues(repo)
	if err != nil {
		// Fallback to manual input
		input := prompt.Input("Issue number", "")
		num, err := strconv.Atoi(input)
		if err != nil {
			return 0, fmt.Errorf("invalid issue number: %s", input)
		}
		return num, nil
	}

	if len(issues) == 0 {
		input := prompt.Input("No issues assigned to you. Enter issue number", "")
		num, err := strconv.Atoi(input)
		if err != nil {
			return 0, fmt.Errorf("invalid issue number: %s", input)
		}
		return num, nil
	}

	options := make([]string, len(issues))
	for i, issue := range issues {
		options[i] = fmt.Sprintf("#%d - %s", issue.Number, issue.Title)
	}

	idx, err := prompt.Select("Select an issue:", options)
	if err != nil {
		return 0, err
	}

	return issues[idx].Number, nil
}

func inferIssueType(issue *ghapi.Issue) string {
	for _, label := range issue.Labels {
		name := label.Name
		switch {
		case contains(name, "bug", "fix"):
			return "bugfix"
		case contains(name, "feature", "enhancement"):
			return "feature"
		case contains(name, "hotfix", "urgent", "critical"):
			return "hotfix"
		case contains(name, "docs", "documentation"):
			return "docs"
		case contains(name, "refactor"):
			return "refactor"
		case contains(name, "test"):
			return "test"
		case contains(name, "chore", "maintenance"):
			return "chore"
		case contains(name, "internal"):
			return "internal"
		}
	}
	return ""
}

func contains(s string, substrs ...string) bool {
	s = toLower(s)
	for _, sub := range substrs {
		if s == sub || len(s) > len(sub) && (s[:len(sub)] == sub || s[len(s)-len(sub):] == sub) {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
