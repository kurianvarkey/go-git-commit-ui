package handlers_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/handlers"
	"github.com/stretchr/testify/assert"
)

// push.go methods
func TestPushSuccess(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			fmt.Println(cmd)
			expected := "git push -u origin main"
			assert.Equal(t, expected, cmd)
			return "pushed successfully", nil
		},
	}

	handlers.PushToOrigin(mock, "main")
}

func TestPushToOriginFailure(t *testing.T) {
	mock := &MockGitHelper{
		ExecuteCommandFunc: func(cmd string) (string, error) {
			expected := "git push -u origin main"
			assert.Equal(t, expected, cmd)
			return "some error output", errors.New("push failed")
		},
	}

	handlers.PushToOrigin(mock, "main")
}
