package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Cli struct{}

func (c *Cli) CurrentBranch() string {
	return c.gitExec("branch", "--show-current")
}

func (c *Cli) CommitsBetween(start string, end string) []string {
	// Create a string in the format "start..end" to pass to git
	commitRange := strings.Join([]string{start, "..", end}, "")

	commits := c.gitExec("show", "--no-patch", "--format=- %h - %s", commitRange)

	lines := strings.Split(commits, "\n")

	return lines
}

func (c *Cli) RootDir() string {
	return c.gitExec("rev-parse", "--show-toplevel")
}

func (c *Cli) NewBranch(branchName string) {
	c.gitExec("switch", "-c", branchName)
}

func (c *Cli) gitExec(args ...string) string {
	command := exec.Command("git", args...)

	out, err := command.CombinedOutput()

	if err != nil {
		// Get error text and display it
		fmt.Println(string(out))

		panic(err)
	}

	return strings.TrimSpace(string(out))
}
