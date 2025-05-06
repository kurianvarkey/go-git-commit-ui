package helpers

type DefaultGitHelper struct{}

func (g *DefaultGitHelper) ExecuteCommand(cmd string) (string, error) {
	return ExecuteCommand(cmd)
}

func (g *DefaultGitHelper) ShowConfirm(message string, defaultYes ...bool) bool {
	return ShowConfirm(message, defaultYes...)
}
