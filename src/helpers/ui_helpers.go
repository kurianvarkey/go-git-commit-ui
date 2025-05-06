package helpers

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/kurianvarkey/gitcommitui/src/settings"
)

var confirmPromptFunc = defaultConfirmPrompt

// ShowSpinner shows a spinner with a given title and executes the given action.
// Spinner type defaults to spinner.Dots if not provided.
func ShowSpinner(title string, action func(), spinnerType ...spinner.Type) {
	var sType spinner.Type = spinner.Dots

	if len(spinnerType) > 0 {
		sType = spinnerType[0]
	}

	if err := spinner.New().Type(sType).Title(title).Action(action).Run(); err != nil {
		return
	}
}

// ShowConfirm displays a confirmation prompt with the specified title.
// The user can confirm or decline the prompt. The default value is used
// if provided, otherwise, it defaults to false.
// Returns true if confirmed, false otherwise.
func ShowConfirm(title string, defaultValue ...bool) bool {
	return confirmPromptFunc(title, defaultValue...)
}

// defaultConfirmPrompt displays a confirmation prompt with the specified title.
// The user can confirm or decline the prompt. The default value is used
// if provided, otherwise, it defaults to false.
// Returns true if confirmed, false otherwise.
func defaultConfirmPrompt(title string, defaultValue ...bool) bool {
	var confirm bool = false
	if len(defaultValue) > 0 {
		confirm = defaultValue[0]
	}

	huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(title).
				Affirmative("Yes").
				Negative("No").
				Value(&confirm),
		),
	).WithTheme(settings.HuhTheme).Run()

	return confirm
}

// SetConfirmPromptFunc sets the function to be used by ShowConfirm to prompt the
// user for confirmation. The default is defaultConfirmPrompt.
func SetConfirmPromptFunc(f func(string, ...bool) bool) {
	confirmPromptFunc = f
}

// ConfirmPromptFunc returns the current confirmation prompt function used by ShowConfirm.
// This function is responsible for displaying a confirmation prompt with the specified title.
// The returned function takes a title and an optional boolean default value, and returns true
// if the user confirms, or false otherwise.
func GetConfirmPromptFunc() func(string, ...bool) bool {
	return confirmPromptFunc
}
