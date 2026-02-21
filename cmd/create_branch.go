package cmd

import (
	"fmt"
	"strconv"

	"github.com/jesusgpo/gh-buddy/internal/branch"
	"github.com/jesusgpo/gh-buddy/internal/ghapi"
	"github.com/jesusgpo/gh-buddy/internal/git"
	"github.com/jesusgpo/gh-buddy/internal/prompt"
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
The branch type can be one of: feature, bugfix, hotfix, release, chore, docs, refactor, test.`,
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
	cmd.Flags().StringVarP(&issueType, "type", "t", "", "branch type (feature, bugfix, hotfix, release, chore, docs, refactor, test)")
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

	fmt.Printf("📋 Issue #%d: %s\n", issue.Number, issue.Title)

	// Determine issue type
	if issueType == "" {
		issueType = inferIssueType(issue)
	}

	if issueType == "" && !useDefaults {
		types := branch.AllIssueTypeStrings()
		idx, err := prompt.Select("Select branch type:", types)
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

	fmt.Printf("🌿 Creating branch: %s (from %s)\n", branchName, baseBranch)

	// Create the branch
	if err := git.CreateBranchFrom(branchName, baseBranch, "origin"); err != nil {
		return err
	}

	fmt.Printf("✅ Branch %q created and checked out successfully!\n", branchName)

	// Ask to push
	shouldPush := useDefaults || prompt.Confirm("Push branch to origin?", true)
	if shouldPush {
		if err := git.PushBranch("origin", branchName); err != nil {
			return err
		}
		fmt.Println("🚀 Branch pushed to origin")
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
