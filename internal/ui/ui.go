package ui

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

// Info prints a blue informational message.
func Info(format string, a ...any) {
	pterm.Info.Println(fmt.Sprintf(format, a...))
}

// Success prints a green success message.
func Success(format string, a ...any) {
	pterm.Success.Println(fmt.Sprintf(format, a...))
}

// Warning prints a yellow warning message.
func Warning(format string, a ...any) {
	pterm.Warning.Println(fmt.Sprintf(format, a...))
}

// Error prints a red error message to stderr.
func Error(format string, a ...any) {
	pterm.Error.Println(fmt.Sprintf(format, a...))
}

// Fatal prints a red error message and exits.
func Fatal(format string, a ...any) {
	pterm.Error.Println(fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Header prints a styled section header.
func Header(title string) {
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		WithFullWidth(false).
		Println(title)
}

// IssuePanel renders a nice box with issue details.
func IssuePanel(number int, title string) {
	content := pterm.Sprintf("  %s %s",
		pterm.FgLightCyan.Sprint(fmt.Sprintf("#%d", number)),
		pterm.FgLightWhite.Sprint(title),
	)
	pterm.DefaultBox.
		WithTitle(pterm.FgYellow.Sprint("Issue")).
		WithTitleTopLeft().
		Println(content)
}

// BranchPanel renders a panel with branch creation details.
func BranchPanel(branchName, baseBranch string) {
	content := pterm.Sprintf("  %s  %s  %s",
		pterm.FgLightGreen.Sprint(branchName),
		pterm.FgGray.Sprint("from"),
		pterm.FgLightYellow.Sprint(baseBranch),
	)
	pterm.DefaultBox.
		WithTitle(pterm.FgYellow.Sprint("New Branch")).
		WithTitleTopLeft().
		Println(content)
}

// PRSummaryPanel renders a summary panel before creating a PR.
func PRSummaryPanel(title, fromBranch, toBranch string, draft bool) {
	rows := [][]string{
		{pterm.FgLightYellow.Sprint("Title"), title},
		{pterm.FgLightYellow.Sprint("From"), pterm.FgLightCyan.Sprint(fromBranch)},
		{pterm.FgLightYellow.Sprint("Into"), pterm.FgLightGreen.Sprint(toBranch)},
	}
	if draft {
		rows = append(rows, []string{pterm.FgLightYellow.Sprint("Draft"), pterm.FgLightMagenta.Sprint("yes")})
	}

	tableStr, _ := pterm.DefaultTable.
		WithHasHeader(false).
		WithData(rows).
		Srender()

	pterm.DefaultBox.
		WithTitle(pterm.FgYellow.Sprint("Pull Request")).
		WithTitleTopLeft().
		Println(tableStr)
}

// BodyPreview renders the PR body in a styled box.
func BodyPreview(body string) {
	pterm.DefaultBox.
		WithTitle(pterm.FgYellow.Sprint("PR Body Preview")).
		WithTitleTopLeft().
		Println(body)
}

// StartSpinner starts a spinner with a label and returns it (call Stop() when done).
func StartSpinner(label string) (*pterm.SpinnerPrinter, error) {
	return pterm.DefaultSpinner.
		WithStyle(pterm.NewStyle(pterm.FgLightCyan)).
		Start(label)
}
