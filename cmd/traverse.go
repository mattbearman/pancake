package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Move up one layer and check out its code",
	Long:  `Move up one layer in the current stack and check out its code`,
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCli := &git.Cli{}

		stackList := *stacks.LoadList(gitCli)

		return up(stackList, gitCli, os.Stdout)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Move down one layer and check out its code",
	Long:  `Move down one layer in the current stack and check out its code`,
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCli := &git.Cli{}

		stackList := *stacks.LoadList(gitCli)

		return down(stackList, gitCli, os.Stdout)
	},
}

type GitForTraverseCmd interface {
	CurrentBranch() string
	ChangeBranch(branchName string)
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}

func up(l stacks.List, g GitForTraverseCmd, out io.Writer) error {
	currentBranch := g.CurrentBranch()
	stack, err := l.ForBranch(currentBranch)

	if err != nil {
		return err
	}

	currentLayer := stack.CurrentLayer()

	newLayer := stack.UpLayer()

	if newLayer == currentLayer {
		return fmt.Errorf("❌ you're already at the top of this stack")
	}

	branchName := stack.BranchForLayer(newLayer)

	g.ChangeBranch(branchName)

	fmt.Fprintf(out, "🥞 Working in git branch \"%s\"\n", branchName)

	return nil
}

func down(l stacks.List, g GitForTraverseCmd, out io.Writer) error {
	currentBranch := g.CurrentBranch()
	stack, err := l.ForBranch(currentBranch)

	if err != nil {
		return err
	}

	currentLayer := stack.CurrentLayer()

	newLayer := stack.DownLayer()

	if newLayer == currentLayer {
		return fmt.Errorf("❌ you're already at the bottom of this stack")
	}

	branchName := stack.BranchForLayer(newLayer)

	g.ChangeBranch(branchName)

	fmt.Fprintf(out, "🥞 Working in git branch \"%s\"\n", branchName)

	return nil
}
