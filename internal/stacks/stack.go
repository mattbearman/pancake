package stacks

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattbearman/pancake/internal/git"
)

type stack struct {
	Name, BaseBranch, Separator string
	Layers                      []string
}

func (s *stack) AddLayer(unscopedBranchName string) (fullBranchName string) {
	fullBranchName = strings.Join([]string{s.Name, s.Separator, unscopedBranchName}, "")

	s.Layers = append(s.Layers, fullBranchName)

	return
}

func (s *stack) UpLayer() (branchName *string) {
	currentLayer := git.CurrentBranch()
	currentIndex := s.LayerIndex(currentLayer)
	layerCount := len(s.Layers)
	nextIndex := currentIndex + 1

	if nextIndex == layerCount {
		return nil
	}

	return &s.Layers[nextIndex]
}

func (s *stack) DownLayer() (branchName *string) {
	currentLayer := git.CurrentBranch()
	currentIndex := s.LayerIndex(currentLayer)

	if currentIndex == 0 {
		return nil
	}

	return &s.Layers[currentIndex-1]
}

func (s *stack) LayerIndex(desiredLayer string) int {
	layerIndex := -1

	s.EachLayer(func(index int, layer string) {
		if layer == desiredLayer {
			layerIndex = index
		}
	})

	if layerIndex == -1 {
		fmt.Printf("❌ branch %s is not part of the %s stack\n", desiredLayer, s.Name)
		os.Exit(1)
	}

	return layerIndex
}

func (s *stack) EachLayer(iterator func(int, string)) {
	for i := 0; i < len(s.Layers); i++ {
		iterator(i, s.Layers[i])
	}
}
