package handlers_test

import (
	"errors"
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/handlers"
	"github.com/kurianvarkey/gitcommitui/src/settings"
	"github.com/stretchr/testify/assert"
)

type MockForm struct {
	RunFunc              func() error
	SetDefaultValuesFunc func(commitTypes []string, defaultCommitType string, defaultVersion string, defaultJiraReference string)
	GetValuesFunc        func() (string, string, string, string)
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

func TestShowCommitUISuccess(t *testing.T) {
	form := &MockForm{
		RunFunc: func() error { return nil },
		GetValuesFunc: func() (string, string, string, string) {
			return "1.0", "feat", "JIRA-123", "Initial commit"
		},
	}

	helper := &MockGitHelper{
		ShowConfirmFunc: func(msg string, defaultYes ...bool) bool {
			assert.Contains(t, msg, "Initial commit")
			return true
		},
		ExecuteCommandFunc: func(cmd string) (string, error) {
			assert.Contains(t, cmd, "Initial commit")
			return "Committed", nil
		},
	}

	config := &settings.Config{
		CommitFormat: "$version-$type-$jira-$summary",
	}

	confirmed := handlers.ShowCommitUI(helper, config, form)
	assert.True(t, confirmed)
}

func TestShowCommitUIFormCancelled(t *testing.T) {
	form := &MockForm{
		RunFunc: func() error { return errors.New("form cancelled") },
	}

	helper := &MockGitHelper{}
	config := &settings.Config{}
	confirmed := handlers.ShowCommitUI(helper, config, form)

	assert.False(t, confirmed)
}
