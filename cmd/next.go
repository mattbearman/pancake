package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

type GitForNextCmd interface {
	CurrentBranch() string
	NewBranch(branchName string)
}

var nextCmd = &cobra.Command{
	Use:   "next [layer name (optional)]",
	Short: "Start a new layer in the current stack",
	Long: `Start a new layer in the current stack, layer name will default to "part-X" if not provided, where X is the the current number of layers +1

	Eg: if your stack has 4 layers currently, calling pancake next will create a new layer called part-5`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gitCli := git.Cli{}

		stackList := *stacks.LoadList(&gitCli)

		exitCode := next(stackList, &gitCli, os.Stdout, args...)

		os.Exit(exitCode)
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}

func next(l stacks.List, g GitForNextCmd, out io.Writer, args ...string) int {
	currentBranch := g.CurrentBranch()
	stack, error := l.ForBranch(currentBranch)

	if error != nil {
		fmt.Fprintln(out, error)

		return 1
	}

	// TODO: Ensure we're currently at the top of the stack

	var layerName string

	if len(args) == 1 {
		layerName = args[0]
	} else {
		layerName = fmt.Sprintf("part-%d", len(stack.Layers)+1)
	}

	stack.AddLayer(layerName)

	branchName := stack.BranchForLayer(layerName)
	g.NewBranch(branchName)

	l.Save()

	fmt.Fprintf(out, "🥞 Created new layer \"%s\". Working in git branch \"%s\"\n", layerName, branchName)

	return 0
}
