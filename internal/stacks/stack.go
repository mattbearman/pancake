package stacks

import (
  "strings"
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
