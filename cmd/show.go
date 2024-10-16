package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display details of the current working stack",
	Long:  `Display details of the current working stack listing all the layers, and the commits in each layer`,
	Run: func(cmd *cobra.Command, args []string) {
		gitCli := git.Cli{}

		stackList := *stacks.LoadList(&gitCli)

		show(stackList, &gitCli, os.Stdout)
	},
}

func show(l stacks.List, g git.Driver, out io.Writer) {
	currentBranch := g.CurrentBranch()
	stack := l.ForBranch(currentBranch)

	fmt.Fprintf(out, "🥞 Current stack: %s\n", stack.Name)
	fmt.Fprintln(out, "   Layers:")

	previousBranch := stack.BaseBranch

	stack.EachLayer(func(_ int, layer string) {
		fmt.Fprintf(out, "   - %s\n", layer)

		currentBranch := stack.BranchForLayer(layer)

		lines := g.CommitsBetween(previousBranch, currentBranch)
		for i := 0; i < len(lines); i++ {
			fmt.Fprintf(out, "     %s\n", lines[i])
		}

		fmt.Fprintln(out)

		previousBranch = currentBranch
	})
}
