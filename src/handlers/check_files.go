package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/kurianvarkey/gitcommitui/src/commands"
	"github.com/kurianvarkey/gitcommitui/src/helpers"
)

// GetStagedFiles gets the list of staged files in the Git repository. If there are
// staged files, it prompts the user to continue with the given files. If the
// user chooses to exit, it returns an empty list of files and 'true' for exit.
//
// If there are no staged files, it returns an empty list of files and 'false' for
// exit.
func GetStagedFiles(helper helpers.GitHelper) (files []string, exit bool) {
	output, err := helper.ExecuteCommand(commands.GitStagedFiles)
	if err != nil {
		return nil, false
	}

	files = strings.Split(strings.TrimSpace(output), "\n")
	if len(files) == 1 && files[0] == "" {
		return nil, false
	}

	if !helper.ShowConfirm(fmt.Sprintf("You have %d staged files. Do you want to continue with following files?\n-> %s", len(files), strings.Join(files, "\n-> ")), true) {
		return []string{}, true
	}

	return files, false
}

// GetChangedFiles gets the list of changed files in the Git repository. If there are
// changes, it prompts the user to stage and continue with the given files. If the
// user chooses to exit, it returns an empty list of files and 'true' for exit.
//
// If there are no changed files, it returns an empty list of files and 'false' for
// exit.
func GetChangedFiles(helper helpers.GitHelper) (files []string, exit bool) {
	output, err := helper.ExecuteCommand(commands.GitChangedFiles)
	if err != nil {
		return nil, false
	}

	files = strings.Split(strings.TrimSpace(output), "\n")
	if len(files) == 1 && files[0] == "" {
		return nil, false
	}

	var changedFiles []string
	for _, file := range files {
		if len(file) >= 3 { // Porcelain output typically follows "XY filename"
			filePath := strings.TrimSpace(file[2:]) // Extract the file path
			if filePath != "" {
				changedFiles = append(changedFiles, filePath)
			}
		}
	}

	if len(changedFiles) > 0 {
		if !helper.ShowConfirm(fmt.Sprintf("You have %d changed files. Do you want to stage and continue with following files?\n-> %s", len(changedFiles), strings.Join(changedFiles, "\n-> ")), true) {
			return []string{}, true
		}

		_, err := helper.ExecuteCommand(commands.GitCommitAdd)
		if err != nil {
			log.Printf("Error staging files: %v", err)
			return nil, false
		}
	}

	return changedFiles, false
}
