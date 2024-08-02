package cmd

import (
	"fmt"
	"os"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Move up one layer and check out its code",
	Long:  `Move up one layer in the current stack and check out its code`,
	Run: func(cmd *cobra.Command, args []string) {
		stacks.Load()

		stack := stacks.Current()

		nextLayer := stack.UpLayer()

		if nextLayer != nil {
			git.ChangeBranch(*nextLayer)
		} else {
			fmt.Println("❌ you're already at the top of this stack")
			os.Exit(1)
		}
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Move down one layer and check out its code",
	Long:  `Move down one layer in the current stack and check out its code`,
	Run: func(cmd *cobra.Command, args []string) {
		stacks.Load()

		stack := stacks.Current()

		nextLayer := stack.DownLayer()

		if nextLayer != nil {
			git.ChangeBranch(*nextLayer)
		} else {
			fmt.Println("❌ you're already at the bottom of this stack")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}
