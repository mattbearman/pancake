package cmd

import (
	"fmt"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

var nextCmd = &cobra.Command{
	Use:   "next [branch name (optional)]",
	Short: "Start a new layer in the current stack",
	Long: `Start a new layer in the current stack, branch name will default to "part-X" if not provided, where X is the the current number of layers +1

	Eg: if your stack has 4 layers currently, calling pancake next will create a new layer called part-5`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stacks.Load()

		stack := stacks.Current()

		// TODO: Ensure we're currently at the top of the stack

		var branchName string

		if len(args) == 1 {
			branchName = args[0]
		} else {
			branchName = fmt.Sprintf("part-%d", len(stack.Layers)+1)
		}

		layer := stack.AddLayer(branchName)

		git.NewBranch(layer)

		stacks.Save()
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
}
