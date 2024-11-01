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
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCli := git.Cli{}

		stackList := *stacks.LoadList(&gitCli)

		return next(stackList, &gitCli, os.Stdout, args...)
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}

func next(l stacks.List, g GitForNextCmd, out io.Writer, args ...string) error {
	currentBranch := g.CurrentBranch()
	stack, err := l.ForBranch(currentBranch)

	if err != nil {
		return err
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

	return nil
}
