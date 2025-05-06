package helpers_test

import (
	"os/exec"
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/helpers"
	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		hasError bool
	}{
		{"simple command", "git status", []string{"git", "status"}, false},
		{"quoted arg", "git commit -m \"Initial commit\"", []string{"git", "commit", "-m", "Initial commit"}, false},
		{"single quotes", "echo 'hello world'", []string{"echo", "hello world"}, false},
		{"mixed quotes", "sh -c 'echo \"hello\"'", []string{"sh", "-c", "echo \"hello\""}, false},
		{"unclosed quote", "git commit -m 'Initial", []string{"git", "commit", "-m", "Initial"}, false}, // doesn't fail but gives partial args
		{"empty input", "", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := helpers.ParseCommand(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, args)
			}
		})
	}
}

func TestExecuteCommandSuccess(t *testing.T) {
	output, err := helpers.ExecuteCommand("ls")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	t.Logf("Output: %v", output)
}

func TestExecuteCommandFailure(t *testing.T) {
	_, err := helpers.ExecuteCommand("nonexistent_command_xyz")
	assert.Error(t, err)
}

func TestClearTerminal(t *testing.T) {
	// Save original and restore after test
	originalExec := helpers.GetExecCommand()
	defer func() { helpers.SetExecCommand(originalExec) }()

	called := false

	// Mock command
	helpers.SetExecCommand(func(name string, args ...string) *exec.Cmd {
		called = true
		return exec.Command("echo", "cleared")
	})

	helpers.ClearTerminal()

	if !called {
		t.Error("expected ExecCommand to be called")
	}
}
