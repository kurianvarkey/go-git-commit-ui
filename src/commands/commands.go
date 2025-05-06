package commands

const (
	GitCheck         = "git rev-parse --is-inside-work-tree"
	GitInit          = "git init"
	GitChangedFiles  = "git status --untracked-files=all --porcelain"
	GitCommitAdd     = "git add -A"
	GitCommitMessage = "git commit -m '%s'"
	GitCurrentBranch = "git rev-parse --abbrev-ref HEAD"
	GitGetRemote     = "git remote get-url origin"
	GitSetRemote     = "git remote set-url origin %s"
	GitPush          = "git push -u origin %s"
	GitStagedFiles   = "git diff --cached --name-only"

	StagedFilesByExtension = `echo %s | grep '\.%s$'`
	BinExits               = "command -v %s >/dev/null 2>&1"
)
