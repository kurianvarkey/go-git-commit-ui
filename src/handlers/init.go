package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/kurianvarkey/gitcommitui/src/commands"
	"github.com/kurianvarkey/gitcommitui/src/helpers"
)

// CheckForGitInitialise checks whether the current directory is a Git repository, and if
// not, prompts the user to initialise a repository and continues.
func CheckForGitInitialise(helper helpers.GitHelper) (exit bool) {
	if !isGitInitialised(helper) {
		if !helper.ShowConfirm("Git repository not found. Do you want to initialise and continue?", true) {
			return true
		}
		initialiseGit(helper)
	}
	return false
}

// IsGitInitialised checks whether the current directory is a Git repository.
func isGitInitialised(helper helpers.GitHelper) bool {
	output, err := helper.ExecuteCommand(commands.GitCheck)
	if err != nil {
		return false
	}
	return strings.TrimSpace(output) == "true"
}

// GetCurrentBranch retrieves the current branch name from the Git repository.
// It executes the command 'git rev-parse --abbrev-ref HEAD' and returns the
// branch name as a string, or an error if the command fails.
func GetCurrentBranch(helper helpers.GitHelper) (string, error) {
	branchName, err := helper.ExecuteCommand(commands.GitCurrentBranch)
	if err != nil {
		return "", fmt.Errorf("failed to determine current branch: %w", err)
	}
	return strings.TrimSpace(branchName), nil
}

// GetRemoteURL retrieves the URL of the remote repository configured for the
// current Git repository. It executes the command 'git remote get-url origin' and
// returns the URL as a string, or an error if the command fails.
func GetRemoteURL(helper helpers.GitHelper) (string, error) {
	url, err := helper.ExecuteCommand(commands.GitGetRemote)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(url), nil
}

// Initialise a Git repository in the current directory.
func initialiseGit(helper helpers.GitHelper) {
	_, err := helper.ExecuteCommand(commands.GitInit)
	if err != nil {
		log.Printf("Failed to initialise git repository: %v", err)
	}
}
