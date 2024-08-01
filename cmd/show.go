/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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
			fmt.Printf("❌🥞 branch %s is not part of a stack\n", currentBranch)

			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
