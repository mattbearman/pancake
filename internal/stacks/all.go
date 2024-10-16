package stacks

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattbearman/pancake/internal/git"
)

var allStacks []*Stack

func Current() *Stack {
	return ForBranch(git.CurrentBranch())
}

func ForBranch(branchName string) *Stack {
	for s := 0; s < len(allStacks); s++ {
		stack := allStacks[s]

		if strings.HasPrefix(branchName, stack.Name) {
			for l := 0; l < len(stack.Layers); l++ {
				if stack.Layers[l] == branchName {
					return stack
				}
			}
		}
	}

	fmt.Printf("❌ branch %s is not part of a stack\n", branchName)
	os.Exit(1)

	return nil
}

func Add(name string, baseBranch string, separator string) *Stack {
	s := Stack{Name: name, BaseBranch: baseBranch, Separator: separator}

	allStacks = append(allStacks, &s)

	return &s
}
