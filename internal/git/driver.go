package git

type Driver interface {
	CurrentBranch() string
	CommitsBetween(start string, end string) []string
	RootDir() string
	NewBranch(branchName string)
}
