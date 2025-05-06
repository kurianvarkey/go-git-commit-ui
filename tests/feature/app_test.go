package feature_test

import (
	"errors"
	"os"
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/cmd"
	"github.com/stretchr/testify/require"
)

type MockGitHelper struct {
	IsRepo       bool
	ChangedFiles []string
	StagedFiles  []string
	BranchName   string
	RemoteURL    string
	ShouldPush   bool
}

type MockForm struct {
	RunFunc              func() error
	SetDefaultValuesFunc func(commitTypes []string, defaultCommitType string, defaultVersion string, defaultJiraReference string)
	GetValuesFunc        func() (string, string, string, string)
}

func (m *MockGitHelper) IsGitRepository() bool {
	return m.IsRepo
}

func (m *MockGitHelper) InitialiseRepository() error {
	return nil
}

func (m *MockGitHelper) GetChangedFiles() ([]string, error) {
	return m.ChangedFiles, nil
}

func (m *MockGitHelper) StageAllChanges() error {
	return nil
}

func (m *MockGitHelper) GetStagedFiles() ([]string, error) {
	return m.StagedFiles, nil
}

func (m *MockGitHelper) Commit(message string) error {
	return nil
}

func (m *MockGitHelper) GetCurrentBranch() (string, error) {
	return m.BranchName, nil
}

func (m *MockGitHelper) Push(branch string) error {
	if m.ShouldPush {
		return nil
	}
	return errors.New("push failed")
}

func (m *MockGitHelper) GetRemoteURL() (string, error) {
	return m.RemoteURL, nil
}

func (m *MockGitHelper) SetRemoteURL(url string) error {
	return nil
}

func (m *MockGitHelper) ShowConfirm(title string, defaultValue ...bool) bool {
	// Return a predefined mock value or just always return true for simplicity
	return true
}

func (m *MockGitHelper) ExecuteCommand(cmd string) (string, error) {
	return "mocked output", nil
}

func (m *MockForm) Run() error {
	if m.RunFunc != nil {
		return m.RunFunc()
	}
	return nil
}

func (m *MockForm) SetDefaultValues(commitTypes []string, defaultCommitType string, defaultVersion string, defaultJiraReference string) {
	if m.SetDefaultValuesFunc != nil {
		m.SetDefaultValuesFunc(commitTypes, defaultCommitType, defaultVersion, defaultJiraReference)
	}
}

func (m *MockForm) GetValues() (string, string, string, string) {
	if m.GetValuesFunc != nil {
		return m.GetValuesFunc()
	}
	return "1.0", "feat", "JIRA-123", "Initial commit"
}

const testConfigFile = "config.json"

func cleanupConfigFile(t *testing.T) {
	t.Helper()
	if _, err := os.Stat(testConfigFile); err == nil {
		if err := os.Remove(testConfigFile); err != nil {
			t.Fatalf("Failed to clean up config file: %v", err)
		}
	}
}

// TestRunAppSuccessfulFlow tests the happy path of the RunApp function by
// mocking all the underlying methods to return successful values. It tests
// that the function returns no error when all the underlying calls are
// successful.
func TestFeatureRunAppSuccessfulFlow(t *testing.T) {
	defer cleanupConfigFile(t)

	mock := &MockGitHelper{
		IsRepo:       true,
		ChangedFiles: []string{"file.txt"},
		BranchName:   "main",
		RemoteURL:    "git@github.com:user/repo.git",
		ShouldPush:   true,
	}

	err := cmd.RunApp(mock, &MockForm{})
	require.NoError(t, err)
}
