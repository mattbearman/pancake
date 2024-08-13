package cmd

import (
	"fmt"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display details of the current working stack",
	Long:  `Display details of the current working stack listing all the layers, and the commits in each layer`,
	Run: func(cmd *cobra.Command, args []string) {
		stacks.Load()

		currentBranch := git.CurrentBranch()
		stack := stacks.ForBranch(currentBranch)

		fmt.Printf("🥞 Current stack: %s\n", stack.Name)
		fmt.Println("   Layers:")

		previousLayer := stack.BaseBranch

		stack.EachLayer(func(_ int, layer string) {
			fmt.Printf("   - %s\n", layer)
			fmt.Println(git.CommitsBetween(previousLayer, layer))
			previousLayer = layer
		})
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
