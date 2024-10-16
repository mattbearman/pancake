package git

import (
	"fmt"
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

func ChangeBranch(branchName string) {
	gitExec("switch", branchName)
}

func CommitsBetween(start string, end string) string {
	// Create a string in the format "start..end" to pass to git
	commitRange := strings.Join([]string{start, "..", end}, "")

	commits := gitExec("show", "--no-patch", "--format=- %h - %s", commitRange)

	lines := strings.Split(commits, "\n")

	for i := 0; i < len(lines); i++ {
		// indent each line as doing it in the --format flag doesn't work
		lines[i] = fmt.Sprintf("     %s", lines[i])
	}

	return strings.Join(lines, "\n")
}

func gitExec(args ...string) string {
	command := exec.Command("git", args...)

	out, err := command.Output()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(out))
}
