package stacks

import (
	"strings"
)

type Stack struct {
	Name              string   `json:"name"`
	BaseBranch        string   `json:"baseBranch"`
	Separator         string   `json:"separator"`
	Layers            []string `json:"layers"`
	currentLayerIndex int
}

func (s *Stack) AddLayer(layerName string) {
	s.Layers = append(s.Layers, layerName)
}

func (s *Stack) UpLayer() string {
	nextIndex := s.currentLayerIndex + 1

	if nextIndex >= len(s.Layers) {
		return s.Layers[s.currentLayerIndex]
	}

	return s.Layers[nextIndex]
}

func (s *Stack) DownLayer() string {
	prevIndex := s.currentLayerIndex - 1

	if prevIndex < 0 {
		return s.Layers[s.currentLayerIndex]
	}

	return s.Layers[prevIndex]
}

func (s *Stack) CurrentLayer() string {
	return s.Layers[s.currentLayerIndex]
}

func (s *Stack) EachLayer(iterator func(int, string)) {
	for i := 0; i < len(s.Layers); i++ {
		iterator(i, s.Layers[i])
	}
}

func (s *Stack) BranchForLayer(layer string) string {
	return strings.Join([]string{s.Name, s.Separator, layer}, "")
}
