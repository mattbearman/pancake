package stacks

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadList(g GitInterface) *List {
	pancakeFile := filepath.Join(g.RootDir(), ".pancake.json")

	dat, err := os.ReadFile(pancakeFile)

	list := NewList(pancakeFile)

	if err != nil {
		// Swallow file not found error, as we will initialize an empty slice of stacks
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	} else {
		json.Unmarshal([]byte(dat), &list.Stacks)
	}

	return list
}

type GitInterface interface {
	RootDir() string
}

func NewList(file string) *List {
	return &List{file: file}
}

type List struct {
	Stacks []*Stack
	file   string
}

func (l *List) ForBranch(branchName string) (*Stack, error) {
	for s := 0; s < len(l.Stacks); s++ {
		stack := l.Stacks[s]

		if strings.HasPrefix(branchName, stack.Name) {
			for l := 0; l < len(stack.Layers); l++ {
				layerBranch := stack.BranchForLayer(stack.Layers[l])
				if layerBranch == branchName {
					return stack, nil
				}
			}

			// Remove the stack name and separator from the branch name to get the name of the layer we couldn't find
			expectedLayer, _ := strings.CutPrefix(branchName, fmt.Sprintf("%s%s", stack.Name, stack.Separator))

			return nil, fmt.Errorf(`❌ layer "%s" is not part of the "%s" stack`, expectedLayer, stack.Name)
		}
	}

	return nil, fmt.Errorf(`❌ branch "%s" is not part of a stack`, branchName)
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
