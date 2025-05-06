package handlers_test

import (
	"errors"
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/commands"
	"github.com/kurianvarkey/gitcommitui/src/handlers"
	"github.com/stretchr/testify/assert"
)

// init.go methods
func TestCheckForGitInitialiseAlreadyInitialised(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "true", nil
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return true
		},
	}

	exit := handlers.CheckForGitInitialise(mock)
	assert.False(t, exit)
}

func TestCheckForGitInitialiseUserDeclinesInit(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "false", nil
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return false
		},
	}

	exit := handlers.CheckForGitInitialise(mock)
	assert.True(t, exit)
}

func TestCheckForGitInitialiseUserAcceptsInit(t *testing.T) {
	calls := []string{}
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			calls = append(calls, cmd)
			if cmd == commands.GitCheck {
				return "false", nil
			}
			if cmd == commands.GitInit {
				return "", nil
			}
			return "", nil
		},
		ShowConfirmFunc: func(message string, defaultYes ...bool) bool {
			return true
		},
	}

	exit := handlers.CheckForGitInitialise(mock)
	assert.False(t, exit)
	assert.Contains(t, calls, commands.GitInit)
}

func TestGetCurrentBranchSuccess(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "main", nil
		},
	}

	branch, err := handlers.GetCurrentBranch(mock)
	assert.NoError(t, err)
	assert.Equal(t, "main", branch)
}

func TestGetCurrentBranchError(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "", errors.New("git error")
		},
	}

	branch, err := handlers.GetCurrentBranch(mock)
	assert.Error(t, err)
	assert.Empty(t, branch)
}

func TestGetRemoteURLSuccess(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "https://github.com/user/repo.git", nil
		},
	}

	url, err := handlers.GetRemoteURL(mock)
	assert.NoError(t, err)
	assert.Equal(t, "https://github.com/user/repo.git", url)
}

func TestGetRemoteURLError(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			return "", errors.New("fetch failed")
		},
	}

	url, err := handlers.GetRemoteURL(mock)
	assert.Error(t, err)
	assert.Empty(t, url)
}
