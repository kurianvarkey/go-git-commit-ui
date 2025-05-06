package handlers

import (
	"fmt"
	"log"

	"github.com/kurianvarkey/gitcommitui/src/commands"
	"github.com/kurianvarkey/gitcommitui/src/helpers"
)

// pushToOrigin executes the Git command to push the current branch to the remote
// repository on the branch with the given name.
//
// If the command fails, an error is logged with the command and output.
func PushToOrigin(helper helpers.GitHelper, remoteBranch string) error {
	cmd := fmt.Sprintf(commands.GitPush, remoteBranch)
	output, err := helper.ExecuteCommand(cmd)
	if err != nil {
		log.Printf("Failed to push to origin: %v\nCommand: %q\nOutput: %q", err, cmd, output)
		return err
	}

	return nil
}
