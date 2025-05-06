package handlers_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/commands"
	"github.com/kurianvarkey/gitcommitui/src/handlers"
	"github.com/stretchr/testify/assert"
)

type MockGitHelper struct {
	ExecuteCommandFunc func(cmd string) (string, error)
	ShowConfirmFunc    func(message string, defaultYes ...bool) bool
}

func (m *MockGitHelper) ExecuteCommand(cmd string) (string, error) {
	if m.ExecuteCommandFunc != nil {
		return m.ExecuteCommandFunc(cmd)
	}
	return "", nil
}

func (m *MockGitHelper) ShowConfirm(message string, defaultYes ...bool) bool {
	if m.ShowConfirmFunc != nil {
		return m.ShowConfirmFunc(message, defaultYes...)
	}

	if len(defaultYes) > 0 {
		return defaultYes[0]
	}

	return true
}

// check_files.go methods

// GetStagedFile method test
func TestGetStagedFilesConfirmYes(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "file1.txt\nfile2.txt", nil
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return true
		},
	}

	files, exit := handlers.GetStagedFiles(mock)
	assert.False(t, exit)
	assert.Equal(t, []string{"file1.txt", "file2.txt"}, files)
}

func TestGetStagedFilesConfirmNo(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "file1.txt\nfile2.txt", nil
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return false
		},
	}

	files, exit := handlers.GetStagedFiles(mock)
	assert.True(t, exit)
	assert.Empty(t, files)
}

func TestGetStagedFilesError(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "", errors.New("git error")
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return true
		},
	}

	files, exit := handlers.GetStagedFiles(mock)
	assert.False(t, exit)
	assert.Empty(t, files)
}

// GetChangedFiles method test
func TestGetChangedFilesConfirmYes(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			if cmd == commands.GitChangedFiles {
				return " M file1.go\n M file2.go", nil
			}
			if cmd == commands.GitCommitAdd {
				return "", nil // simulate successful staging
			}
			return "", fmt.Errorf("unexpected command")
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return true
		},
	}

	files, exit := handlers.GetChangedFiles(mock)
	assert.False(t, exit)
	assert.Equal(t, []string{"file1.go", "file2.go"}, files)
}

func TestGetChangedFilesConfirmNo(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return " M file1.go\n M file2.go", nil
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return false
		},
	}

	files, exit := handlers.GetChangedFiles(mock)
	assert.True(t, exit)
	assert.Empty(t, files)
}

func TestGetChangedFilesNoChanges(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "", nil
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return true
		},
	}

	files, exit := handlers.GetChangedFiles(mock)
	assert.False(t, exit)
	assert.Empty(t, files)
}

func TestGetChangedFilesCommandError(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "", errors.New("git error")
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return true
		},
	}

	files, exit := handlers.GetChangedFiles(mock)
	assert.False(t, exit)
	assert.Empty(t, files)
}
