package main

import (
	"log"
	"os"
	"time"

	"github.com/kurianvarkey/gitcommitui/src/cmd"
	"github.com/kurianvarkey/gitcommitui/src/handlers"
	"github.com/kurianvarkey/gitcommitui/src/helpers"
)

// main is the entry point of the application.
//
// It initialises the application, then runs the main loop which involves
// checking for changed files, staging changes, prompting the user for a
// commit message, committing the changes, and prompting the user to push
// the branch to origin.
//
// If any step fails, it logs the error and exits with a non-zero status code.
func main() {
	log.SetFlags(0)

	helpers.ShowSpinner("Initialising...", func() {
		time.Sleep(1 * time.Second)
	})

	err := cmd.RunApp(&helpers.DefaultGitHelper{}, &handlers.DefaultCommitForm{})
	if err != nil {
		log.Println("Exiting application:", err)
	}

	os.Exit(0)
}
