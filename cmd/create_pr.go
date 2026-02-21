package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jesusgpo/gh-buddy/internal/ghapi"
	"github.com/jesusgpo/gh-buddy/internal/git"
	"github.com/jesusgpo/gh-buddy/internal/prompt"
	"github.com/spf13/cobra"
)

func newCreatePRCmd() *cobra.Command {
	var (
		issueNumber int
		baseBranch  string
		title       string
		body        string
		draft       bool
		labels      []string
	)

	cmd := &cobra.Command{
		Use:   "create-pr",
		Short: "Create a pull request from the current local branch",
		Long: `Create a pull request from the current branch.

If an issue number is detected from the branch name or provided explicitly, the PR 
title and body will be pre-populated from the issue. Supports linking issues 
automatically via "Closes #N" in the PR body.`,
		Example: `  # Create a PR from the current branch (auto-detect issue)
  gh buddy create-pr

  # Create a PR linked to a specific issue
  gh buddy create-pr --issue 42

  # Create a draft PR
  gh buddy create-pr --draft

  # Create a PR with a custom base branch
  gh buddy create-pr --base develop

  # Use defaults without prompts
  gh buddy create-pr -y`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreatePR(issueNumber, baseBranch, title, body, draft, labels)
		},
	}

	cmd.Flags().IntVarP(&issueNumber, "issue", "i", 0, "issue number to link the PR to")
	cmd.Flags().StringVarP(&baseBranch, "base", "b", "", "base branch for the PR (default: repo default branch)")
	cmd.Flags().StringVarP(&title, "title", "T", "", "PR title (default: generated from issue or branch)")
	cmd.Flags().StringVar(&body, "body", "", "PR body")
	cmd.Flags().BoolVarP(&draft, "draft", "d", false, "create as a draft PR")
	cmd.Flags().StringSliceVarP(&labels, "label", "l", nil, "labels to add to the PR")

	return cmd
}

func runCreatePR(issueNumber int, baseBranch, title, body string, draft bool, labels []string) error {
	repo, err := git.RepoSlug()
	if err != nil {
		return fmt.Errorf("not in a git repository or no origin remote: %w", err)
	}

	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	fmt.Printf("🌿 Current branch: %s\n", currentBranch)

	// Try to detect issue number from branch name
	if issueNumber == 0 {
		issueNumber = extractIssueFromBranch(currentBranch)
	}

	// Fetch issue details if we have a number
	var issue *ghapi.Issue
	if issueNumber > 0 {
		issue, err = ghapi.GetIssue(repo, issueNumber)
		if err != nil {
			fmt.Printf("⚠️  Could not fetch issue #%d: %v\n", issueNumber, err)
		} else {
			fmt.Printf("📋 Linked issue #%d: %s\n", issue.Number, issue.Title)
		}
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

	// Generate title
	if title == "" {
		if issue != nil {
			title = issue.Title
		} else {
			title = generateTitleFromBranch(currentBranch)
		}
		if !useDefaults {
			title = prompt.Input("PR title", title)
		}
	}

	// Generate body
	if body == "" {
		body = generatePRBody(issue)
		if !useDefaults {
			fmt.Println("\n--- PR body preview ---")
			fmt.Println(body)
			fmt.Println("--- end preview ---\n")
			if !prompt.Confirm("Use this PR body?", true) {
				body = prompt.Input("PR body", "")
			}
		}
	}

	// Draft
	if !useDefaults && !draft {
		draft = prompt.Confirm("Create as draft?", false)
	}

	fmt.Printf("\n📝 Creating PR: %s\n", title)
	fmt.Printf("   %s → %s\n", currentBranch, baseBranch)
	if draft {
		fmt.Println("   📌 Draft PR")
	}

	if !useDefaults {
		if !prompt.Confirm("Proceed?", true) {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Push the branch first
	fmt.Println("🚀 Pushing branch to origin...")
	if err := git.PushBranch("origin", currentBranch); err != nil {
		// Branch might already be pushed, continue
		fmt.Printf("⚠️  Push warning: %v (continuing anyway)\n", err)
	}

	pr, err := ghapi.CreatePR(repo, title, body, baseBranch, currentBranch, draft, labels)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Pull request created: %s\n", pr.URL)
	return nil
}

var issueNumberRegex = regexp.MustCompile(`/GH-(\d+)-`)

func extractIssueFromBranch(branchName string) int {
	matches := issueNumberRegex.FindStringSubmatch(branchName)
	if len(matches) >= 2 {
		num, err := strconv.Atoi(matches[1])
		if err == nil {
			return num
		}
	}
	return 0
}

func generateTitleFromBranch(branchName string) string {
	// Remove type prefix (e.g., "feature/")
	parts := strings.SplitN(branchName, "/", 2)
	title := branchName
	if len(parts) == 2 {
		title = parts[1]
	}
	// Remove issue number prefix
	title = regexp.MustCompile(`^\d+-`).ReplaceAllString(title, "")
	// Replace hyphens with spaces and capitalize
	title = strings.ReplaceAll(title, "-", " ")
	if len(title) > 0 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}
	return title
}

func generatePRBody(issue *ghapi.Issue) string {
	var sb strings.Builder

	if issue != nil {
		sb.WriteString("## Description\n\n")
		if issue.Body != "" {
			sb.WriteString(issue.Body)
		} else {
			sb.WriteString(fmt.Sprintf("Resolves #%d", issue.Number))
		}
		sb.WriteString("\n\n")
		sb.WriteString(fmt.Sprintf("Closes #%d\n", issue.Number))
	} else {
		sb.WriteString("## Description\n\n")
		sb.WriteString("<!-- Describe your changes here -->\n\n")
		sb.WriteString("## Checklist\n\n")
		sb.WriteString("- [ ] Tests added/updated\n")
		sb.WriteString("- [ ] Documentation updated\n")
		sb.WriteString("- [ ] Code follows project conventions\n")
	}

	return sb.String()
}
