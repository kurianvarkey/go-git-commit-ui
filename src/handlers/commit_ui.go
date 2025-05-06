package handlers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kurianvarkey/gitcommitui/src/commands"
	"github.com/kurianvarkey/gitcommitui/src/helpers"
	"github.com/kurianvarkey/gitcommitui/src/settings"
)

// Interface to abstract the form
type CommitForm interface {
	Run() error
	SetDefaultValues(commitTypes []string, defaultCommitType string, defaultVersion string, defaultJiraReference string)
	GetValues() (version, commitType, jira, summary string)
}

// Struct to encapsulate form values and logic
type DefaultCommitForm struct {
	Version, CommitType, Jira, Summary string
	Types                              []string
}

var defaultCommitTypes = []string{
	"feat", "fix", "refactor", "chore", "revert",
	"db", "docs", "build", "ci", "perf", "style", "test", "wip",
}

// Run prompts the user to input the commit message details and validates the input.
// It is responsible for rendering the form and handling user input.
func (f *DefaultCommitForm) Run() error {
	commitTypes := f.Types
	if len(commitTypes) == 0 {
		commitTypes = defaultCommitTypes
	}

	options := make([]huh.Option[string], len(commitTypes))
	for i, v := range commitTypes {
		options[i] = huh.NewOption(v, v)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Version").Value(&f.Version).Placeholder("1.x"),
			huh.NewSelect[string]().Title("Commit Type").Options(options...).Value(&f.CommitType),
			huh.NewInput().Title("Reference").Value(&f.Jira).Placeholder("Jira ticket if any"),
			huh.NewText().Title("Summary").Value(&f.Summary).Placeholder("Summary of change").Validate(func(s string) error {
				if s == "" {
					return errors.New("summary cannot be empty")
				}
				return nil
			}),
		),
	).WithTheme(settings.HuhTheme).Run()
}

// SetDefaultValues initializes the commit form fields with the provided values.
// It sets the available commit types, the default commit type, version, and Jira reference.
func (f *DefaultCommitForm) SetDefaultValues(commitTypes []string, defaultCommitType string, defaultVersion string, defaultJiraReference string) {
	f.Types = commitTypes
	f.CommitType = defaultCommitType
	f.Version = defaultVersion
	f.Jira = defaultJiraReference
}

// GetValues returns the values of the commit form fields in the order of version, commit type, jira reference, and summary.
func (f *DefaultCommitForm) GetValues() (string, string, string, string) {
	return f.Version, f.CommitType, f.Jira, f.Summary
}

// ShowCommitUI displays a user interface for inputting commit details using the provided form.
// It prompts the user to confirm committing with the generated commit message.
// If confirmed, it executes the git commit command with the formatted message.
// Returns true if the commit was successful or confirmed, false otherwise.
func ShowCommitUI(helper helpers.GitHelper, config *settings.Config, form CommitForm) bool {
	if err := form.Run(); err != nil {
		return false
	}

	version, commitType, jira, summary := form.GetValues()
	commitMessage := formatCommitMessage(config, version, commitType, jira, summary)

	if helper.ShowConfirm("Commit changes with following message?\n" + commitMessage) {
		_, err := helper.ExecuteCommand(fmt.Sprintf(commands.GitCommitMessage, commitMessage))
		if err != nil {
			log.Printf("Failed to commit changes: %v", err)
		}
		return true
	}

	return false
}

// formatCommitMessage takes in the version, commit type, jira reference, and summary as strings and replaces placeholders in the
// git commit format string with the given values, returning the formatted string.
func formatCommitMessage(config *settings.Config, version string, commitType string, jiraReference string, summary string) string {
	formattedMessage := config.CommitFormat

	formattedMessage = strings.ReplaceAll(formattedMessage, "$version", version)
	formattedMessage = strings.ReplaceAll(formattedMessage, "$type", commitType)
	formattedMessage = strings.ReplaceAll(formattedMessage, "$jira", jiraReference)
	formattedMessage = strings.ReplaceAll(formattedMessage, "$summary", summary)

	return formattedMessage
}
