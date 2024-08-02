package cmd

import (
	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

var separator string

var startCmd = &cobra.Command{
	Use:   "start [stack name] [branch name (optional)]",
	Short: "Start a new stack",
	Long: `Creates a new branch based on the stack name, separator and branch name,
and initializes a stack starting with that branch

Stack name must be provided, branch name will default to "part-1" if not provided,
and the default separator is "/"

Eg:

pancake start new-feature

creates a new branch called "new-feature/part-1"

pancake start new-feature setup --separator=_

creates a new branch called "new-feature_setup"`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(2)),
	Run: func(cmd *cobra.Command, args []string) {
		stacks.Load()

		stackName := args[0]

		branchName := "part-1"

		if len(args) == 2 {
			branchName = args[1]
		}

		stack := stacks.Add(stackName, git.CurrentBranch(), separator)
		layer := stack.AddLayer(branchName)

		git.NewBranch(layer)

		stacks.Save()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&separator, "separator", "s", "/", "a string used to separate stack name and feature name when creating git branches")
}
