package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
)

// ExecCommand is a variable for exec.Command, to allow test injection.
var execCommand = exec.Command

// clearTerminal clears the terminal screen.
func ClearTerminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin": // Linux and macOS
		cmd = execCommand("clear")
	case "windows": // Windows
		// On Windows, 'cls' is an internal command of cmd.exe
		cmd = execCommand("cmd", "/c", "cls")
	default: // Unsupported operating systems
		fmt.Println("Warning: Clearing terminal not supported on this OS.")
		return
	}

	// Set the command's output to the standard output so it clears the user's terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr // Also redirect stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		// Log the error, but don't necessarily stop the program
		fmt.Fprintf(os.Stderr, "Error clearing terminal: %v\n", err)
	}
}

// SetExecCommand sets a custom function for creating exec.Cmd instances.
// This is useful for injecting mock commands in tests.
// The provided function should have the same signature as exec.Command.
func SetExecCommand(f func(string, ...string) *exec.Cmd) {
	execCommand = f
}

// GetExecCommand returns the current function used by ExecuteCommand to create
// exec.Cmd instances. This is useful for getting the current mock command
// function in tests.
func GetExecCommand() func(string, ...string) *exec.Cmd {
	return execCommand
}

// ExecuteCommand executes a command and returns the output
func ExecuteCommand(command string) (string, error) {
	args, err := ParseCommand(strings.TrimSpace(command))
	if err != nil {
		return "", fmt.Errorf("failed to parse command string: %w", err)
	}

	cmd := exec.Command(args[0], args[1:]...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w\nOutput:\n%s", err, string(output))
	}

	return string(output), nil
}

// parseCommand takes a command string and parses it into a slice of strings, respecting quotes to allow for arguments with spaces.
// It returns an error if the command string is empty, or if there's an unclosed quote.
//
// The parsing is done according to the following rules:
// - Any non-space characters that are not quotes are treated as literal characters.
// - When a quote character is encountered, all characters following it are treated as literal until the same quote character is encountered again.
// - Spaces between quotes are treated as literal characters.
// - Any arguments that are not quoted are split on spaces.
// - The last argument is added to the slice regardless of whether it ends with a space or not.
func ParseCommand(cmd string) ([]string, error) {
	args := []string{}
	currentArg := bytes.Buffer{}
	inQuotes := false
	quoteChar := rune(0)

	for _, r := range cmd {
		switch {
		case unicode.IsSpace(r) && !inQuotes:
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		case r == '"' || r == '\'':
			if inQuotes && r == quoteChar {
				inQuotes = false
				quoteChar = 0
			} else if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else {
				// If inside quotes but not the matching quote, treat as literal character
				currentArg.WriteRune(r)
			}
		default:
			currentArg.WriteRune(r)
		}
	}

	// Add the last argument if any
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	// Basic validation: ensure there's at least a command
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command string")
	}

	return args, nil
}
