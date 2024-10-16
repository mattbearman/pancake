package stacks

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattbearman/pancake/internal/git"
)

type Stack struct {
	Name, BaseBranch, Separator string
	Layers                      []string
}

func (s *Stack) AddLayer(layerName string) string {
	s.Layers = append(s.Layers, layerName)

	return s.BranchForLayer(layerName)
}

func (s *Stack) UpLayer() (branchName *string) {
	currentLayer := git.CurrentBranch()
	currentIndex := s.LayerIndex(currentLayer)
	layerCount := len(s.Layers)
	nextIndex := currentIndex + 1

	if nextIndex == layerCount {
		return nil
	}

	return &s.Layers[nextIndex]
}

func (s *Stack) DownLayer() (branchName *string) {
	currentLayer := git.CurrentBranch()
	currentIndex := s.LayerIndex(currentLayer)

	if currentIndex == 0 {
		return nil
	}

	return &s.Layers[currentIndex-1]
}

func (s *Stack) LayerIndex(desiredLayer string) int {
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

func (s *Stack) EachLayer(iterator func(int, string)) {
	for i := 0; i < len(s.Layers); i++ {
		iterator(i, s.Layers[i])
	}
}

func (s *Stack) BranchForLayer(layer string) string {
	return strings.Join([]string{s.Name, s.Separator, layer}, "")
}
