package prompt

import (
	"github.com/pterm/pterm"
)

// Confirm asks the user for a yes/no confirmation using an interactive prompt.
func Confirm(message string, defaultYes bool) bool {
	defaultText := "No"
	if defaultYes {
		defaultText = "Yes"
	}
	result, err := pterm.DefaultInteractiveConfirm.
		WithDefaultText(message).
		WithDefaultValue(defaultYes).
		WithConfirmText("Yes").
		WithRejectText("No").
		Show(defaultText)
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
	selected, err := pterm.DefaultInteractiveSelect.
		WithOptions(options).
		WithMaxHeight(10).
		Show(message)
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
