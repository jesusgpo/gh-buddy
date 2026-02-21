package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// Confirm asks the user for a yes/no confirmation.
func Confirm(message string, defaultYes bool) bool {
	suffix := " [y/N]: "
	if defaultYes {
		suffix = " [Y/n]: "
	}
	fmt.Print(message + suffix)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}

// Input asks the user for text input.
func Input(message, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", message, defaultVal)
	} else {
		fmt.Printf("%s: ", message)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

// Select asks the user to select from a list of options. Returns the index.
func Select(message string, options []string) (int, error) {
	fmt.Println(message)
	for i, opt := range options {
		fmt.Printf("  [%d] %s\n", i+1, opt)
	}
	fmt.Print("Choose an option: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(options) {
		return -1, fmt.Errorf("invalid selection: %s", input)
	}
	return idx - 1, nil
}
