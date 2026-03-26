package cmd

import (
	"os"

	"github.com/jesusgpo/gh-buddy/internal/ui"
	"github.com/spf13/cobra"
)

var (
	version     = "dev"
	useDefaults bool
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "buddy",
		Short: "GitHub CLI Buddy Extension - Your friendly PR & branch companion",
		Long: `GitHub CLI Buddy Extension

Buddy helps you create branches and pull requests following
consistent naming conventions, directly from GitHub issues.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
	}

	rootCmd.PersistentFlags().BoolVarP(&useDefaults, "yes", "y", false, "use the default proposed fields")

	rootCmd.AddCommand(newCreateBranchCmd())
	rootCmd.AddCommand(newCreatePRCmd())

	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}
}
