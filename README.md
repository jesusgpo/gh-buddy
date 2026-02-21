# gh-buddy

> Your friendly GitHub CLI extension for branches & pull requests 🤝

A GitHub CLI extension that helps you create branches and pull requests following consistent naming conventions, directly from GitHub issues.

## Installation

```bash
gh extension install jesusgpo/gh-buddy
```

## Usage

```
gh buddy --help

GitHub CLI Buddy Extension

Buddy helps you create branches and pull requests following
consistent naming conventions, directly from GitHub issues.

Usage:
  buddy [command]

Available Commands:
  create-branch Create a local branch from an issue
  create-pr     Create a pull request from the current local branch
  help          Help about any command

Flags:
  -h, --help      help for buddy
  -v, --version   version for buddy
  -y, --yes       use the default proposed fields
```

### Create a branch

```bash
# Interactive: select from your assigned issues
gh buddy create-branch

# From a specific issue
gh buddy create-branch --issue 42

# With a specific type
gh buddy create-branch --issue 42 --type bugfix

# From a different base branch
gh buddy create-branch --issue 42 --base develop

# Non-interactive: use all defaults
gh buddy create-branch --issue 42 -y
```

Branch naming convention: `<type>/GH-<issue-number>-<slugified-title>`

Supported types: `feature`, `bugfix`, `hotfix`, `release`, `chore`, `docs`, `refactor`, `test`

### Create a pull request

```bash
# Auto-detect issue from branch name
gh buddy create-pr

# Link to a specific issue
gh buddy create-pr --issue 42

# Create as a draft
gh buddy create-pr --draft

# Non-interactive
gh buddy create-pr -y
```

The PR body is auto-generated with:
- Issue description (if linked)
- `Closes #N` reference for automatic issue closing
- Checklist template for unlinked PRs

## Development

```bash
# Build
make build

# Install locally
make install

# Run tests
make test

# Build for all platforms
make release
```

## How it works

1. **create-branch**: Fetches issue details from GitHub, generates a branch name following `type/number-title` convention, creates the branch from the base, and optionally pushes it.

2. **create-pr**: Detects the issue number from the current branch name (or prompts), fetches issue details, generates title/body, pushes the branch, and creates the PR via `gh`.

## Requirements

- [GitHub CLI](https://cli.github.com/) (`gh`) installed and authenticated
- Git

## License

MIT
