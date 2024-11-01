package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/spf13/cobra"
)

type GitForShowCmd interface {
	CurrentBranch() string
	CommitsBetween(from string, to string) []string
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display details of the current working stack",
	Long:  "Display details of the current working stack listing all the layers, and the commits in each layer",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCli := &git.Cli{}

		stackList := *stacks.LoadList(gitCli)

		return show(stackList, gitCli, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func show(l stacks.List, g GitForShowCmd, out io.Writer) error {
	currentBranch := g.CurrentBranch()
	stack, err := l.ForBranch(currentBranch)

	if err != nil {
		return err
	}

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

	return nil
}
