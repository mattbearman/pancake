package cmd

import (
	"fmt"
	"os"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display details of the current working stack",
	Long:  `Display details of the current working stack listing all the layers, and the commits in each layer`,
	Run: func(cmd *cobra.Command, args []string) {
		stacks.Load()

		currentBranch := git.CurrentBranch()
		stack := stacks.ForBranch(currentBranch)

		if stack != nil {
			fmt.Printf("🥞 Current stack: %s\n", currentBranch)
			fmt.Println("   Layers:")

			previousLayer := stack.BaseBranch

			for i := 0; i < len(stack.Layers); i++ {
				layer := stack.Layers[i]
				fmt.Printf("   - %s\n", layer)
				fmt.Println(git.CommitsBetween(previousLayer, layer))
				previousLayer = layer
			}
		} else {
			fmt.Printf("❌ branch %s is not part of a stack\n", currentBranch)

			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
