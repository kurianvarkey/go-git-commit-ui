// cmd/runner.go
package cmd

import (
	"fmt"

	"github.com/kurianvarkey/gitcommitui/src/handlers"
	"github.com/kurianvarkey/gitcommitui/src/helpers"
	"github.com/kurianvarkey/gitcommitui/src/settings"
)

// RunApp is the main entrypoint for the application. It takes a Git helper and
// a commit form as arguments and runs the application logic. It loads the
// configuration, checks for changed files, shows the commit user interface,
// and pushes the changes to the remote repository. If the user cancels at
// any point, it returns an error.
func RunApp(gitHelper helpers.GitHelper, form handlers.CommitForm) error {
	config, err := settings.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Step 1: check whether git is initialised
	if handlers.CheckForGitInitialise(gitHelper) {
		return fmt.Errorf("not a git repository")
	}

	// Step 2: check for changed files
	changedFiles, exit := handlers.GetStagedFiles(gitHelper)
	if exit {
		return fmt.Errorf("no staged files")
	}

	if len(changedFiles) == 0 {
		changedFiles, exit = handlers.GetChangedFiles(gitHelper)
		if exit || len(changedFiles) == 0 {
			return fmt.Errorf("no changed files")
		}
	}

	form.SetDefaultValues(config.CommitTypes, config.DefaultCommitType, config.DefaultVersion, config.DefaultJiraReference)

	if !handlers.ShowCommitUI(gitHelper, config, form) {
		return fmt.Errorf("user canceled commit UI")
	}

	branchName, err := handlers.GetCurrentBranch(gitHelper)
	if err != nil {
		return fmt.Errorf("failed to determine current branch: %w", err)
	}

	if !gitHelper.ShowConfirm(fmt.Sprintf("Do you want to push the current branch '%s' to origin?", branchName), true) {
		return fmt.Errorf("user canceled push")
	}

	remoteURL, err := handlers.GetRemoteURL(gitHelper)
	if err != nil {
		return fmt.Errorf("failed to determine remote URL: %w", err)
	}

	fmt.Printf("Pushing to %s\n", remoteURL)
	if err = handlers.PushToOrigin(gitHelper, branchName); err != nil {
		return fmt.Errorf("failed to push to origin: %w", err)
	}

	fmt.Printf("Pushed to %s successfully\n", remoteURL)

	return nil
}
