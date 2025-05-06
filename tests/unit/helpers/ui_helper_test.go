package helpers_test

import (
	"testing"

	"github.com/kurianvarkey/gitcommitui/src/helpers"
)

func TestShowConfirmMocked(t *testing.T) {
	original := helpers.GetConfirmPromptFunc() // getter we'll define shortly
	defer helpers.SetConfirmPromptFunc(original)

	helpers.SetConfirmPromptFunc(func(title string, defaultValue ...bool) bool {
		if title != "Are you sure?" {
			t.Errorf("Expected title to be 'Are you sure?', got %s", title)
		}
		return true
	})

	result := helpers.ShowConfirm("Are you sure?", true)
	if !result {
		t.Errorf("Expected ShowConfirm to return true, got false")
	}
}

func TestShowSpinnerRunsAction(t *testing.T) {
	called := false
	helpers.ShowSpinner("Testing...", func() {
		called = true
	})
	if !called {
		t.Errorf("Expected action to be called")
	}
}
