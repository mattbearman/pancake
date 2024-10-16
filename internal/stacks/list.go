package stacks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattbearman/pancake/internal/git"
)

func LoadList(g git.Driver) *List {
	pancakeFile := filepath.Join(g.RootDir(), ".pancake.json")

	dat, err := os.ReadFile(pancakeFile)

	list := List{file: pancakeFile}

	if err != nil {
		// Swallow file not found error, as we will initialize an empty slice of stacks
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	} else {
		json.Unmarshal([]byte(dat), &list.Stacks)

	}

	return &list
}

type List struct {
	Stacks []*Stack
	file   string
}

func (l *List) ForBranch(branchName string) *Stack {
	for s := 0; s < len(l.Stacks); s++ {
		stack := l.Stacks[s]

		if strings.HasPrefix(branchName, stack.Name) {
			for l := 0; l < len(stack.Layers); l++ {
				layerBranch := stack.BranchForLayer(stack.Layers[l])
				if layerBranch == branchName {
					return stack
				}
			}
		}
	}

	fmt.Printf("❌ branch %s is not part of a stack\n", branchName)
	os.Exit(1)

	return nil
}

func (l *List) Add(name string, baseBranch string, separator string) *Stack {
	s := Stack{Name: name, BaseBranch: baseBranch, Separator: separator}

	l.Stacks = append(l.Stacks, &s)

	return &s
}

func (l *List) Save() {
	dat, err := json.Marshal(l.Stacks)

	if err != nil {
		panic(err)
	}

	os.WriteFile(l.file, dat, 0644)
}
