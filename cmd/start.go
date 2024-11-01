package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

var (
	separator string

	startCmd = &cobra.Command{
		Use:   "start [stack name] [layer name (optional)]",
		Short: "Start a new stack",
		Long: `Creates a new branch based on the stack name, separator and layer name,
  and initializes a stack starting with that layer

  Stack name must be provided, layer name will default to "part-1" if not provided,
  and the default separator is "/"

  Eg:

  pancake start new-feature

  creates a new branch called "new-feature/part-1"

  pancake start new-feature setup --separator=_

  creates a new branch called "new-feature_setup"`,
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(2)),
		Run: func(cmd *cobra.Command, args []string) {
			gitCli := git.Cli{}

			stackList := *stacks.LoadList(&gitCli)

			start(stackList, &gitCli, os.Stdout, args...)
		},
	}
)

type GitForStartCmd interface {
	CurrentBranch() string
	NewBranch(branchName string)
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&separator, "separator", "s", "/", "a string used to separate stack name and feature name when creating git branches")
}

func start(l stacks.List, g GitForStartCmd, out io.Writer, args ...string) {
	stackName := args[0]

	layerName := "part-1"

	if len(args) == 2 {
		layerName = args[1]
	}

	stack := l.Add(stackName, g.CurrentBranch(), separator)
	stack.AddLayer(layerName)

	branchName := stack.BranchForLayer(layerName)
	g.NewBranch(branchName)

	l.Save()

	fmt.Fprintf(out, "🥞 Created new stack \"%s\". Working in git branch \"%s\"\n", stackName, branchName)
}
