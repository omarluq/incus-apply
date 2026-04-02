// Package terminal provides utilities for terminal detection and output control.
package terminal

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// ANSI color codes for terminal output formatting.
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorDim    = "\033[2m"
	ColorReset  = "\033[0m"
)

// IsTerminal reports whether f is connected to a real terminal.
func IsTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// ClearLine erases the current line and moves the cursor up one line.
// No-op when stdout is not a real terminal.
func ClearLine() {
	if IsTerminal(os.Stdout) {
		fmt.Fprint(os.Stdout, "\033[1A\033[2K\r")
	}
}

// ConfirmPrompt prompts the user for confirmation. Returns true if accepted.
func ConfirmPrompt(prompt string) bool {
	fmt.Printf("\n%s? [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	fmt.Println() // Add a newline after the prompt for cleaner output.

	return response == "y" || response == "yes"
}

// Width returns the current terminal width for the given writer when available.
// If the width cannot be detected, fallback is returned.
func Width(w io.Writer, fallback int) int {
	file, ok := w.(*os.File)
	if !ok {
		return fallback
	}
	width, _, err := term.GetSize(int(file.Fd()))
	if err != nil || width <= 0 {
		return fallback
	}
	return width
}
