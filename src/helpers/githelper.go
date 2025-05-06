package helpers

type GitHelper interface {
	ExecuteCommand(cmd string) (string, error)
	ShowConfirm(message string, defaultYes ...bool) bool
}
