package git

import (
	"os/exec"
	"strings"
)

func CurrentBranch() string {
	return gitExec("branch", "--show-current")
}

func WorkingDir() string {
	return gitExec("rev-parse", "--show-toplevel")
}

func NewBranch(branchName string) {
	gitExec("switch", "-c", branchName)
}

func CommitsBetween(start string, end string) string {
	commitRange := strings.Join([]string{start, "..", end}, "")
	return gitExec("show", "--no-patch", `--format="     - %h - %s"`, commitRange)
}

func gitExec(args ...string) string {
	command := exec.Command("git", args...)

	out, err := command.Output()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(out))
}
