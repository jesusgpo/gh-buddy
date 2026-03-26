package prompt

import (
	"github.com/pterm/pterm"
)

// Confirm asks the user for a yes/no confirmation using an interactive prompt.
func Confirm(message string, defaultYes bool) bool {
	result, err := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(defaultYes).
		Show(message)
	if err != nil {
		return defaultYes
	}
	return result
}

// Input asks the user for text input with an optional default value.
func Input(message, defaultVal string) string {
	result, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText(defaultVal).
		WithDefaultValue(defaultVal).
		Show(message)
	if err != nil || result == "" {
		return defaultVal
	}
	return result
}

// Select asks the user to select from a list of options using arrow keys.
// Returns the index of the selected option.
func Select(message string, options []string) (int, error) {
	return SelectWithDefault(message, options, -1)
}

// SelectWithDefault asks the user to select from a list with a pre-selected default (0-based index).
// Returns the index of the selected option.
func SelectWithDefault(message string, options []string, defaultIdx int) (int, error) {
	printer := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithMaxHeight(10)
	if defaultIdx >= 0 && defaultIdx < len(options) {
		printer = printer.WithDefaultOption(options[defaultIdx])
	}
	selected, err := printer.Show(message)
	if err != nil {
		return -1, err
	}
	for i, opt := range options {
		if opt == selected {
			return i, nil
		}
	}
	return -1, nil
}
