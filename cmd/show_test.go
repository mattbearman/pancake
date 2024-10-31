package cmd

import (
	"bytes"
	"testing"

	"github.com/mattbearman/pancake/internal/stacks"
)

type gitMock struct{}

func (g *gitMock) CurrentBranch() string {
	return "feature/layer2"
}

func (g *gitMock) CommitsBetween(from string, to string) []string {
	switch to {
	case "feature/layer1":
		return []string{"abc123 - commit1", "def456 - commit2"}
	case "feature/layer2":
		return []string{"ghi789 - commit3", "jkl012 - commit4"}
	}

	return []string{}
}

func TestShow(t *testing.T) {
	g := gitMock{}

	stackList := stacks.List{
		Stacks: []*stacks.Stack{
			{
				Name:       "feature",
				BaseBranch: "main",
				Separator:  "/",
				Layers:     []string{"layer1", "layer2"},
			},
		},
	}

	buf := &bytes.Buffer{}
	show(stackList, &g, buf)

	expectedOutput := `🥞 Current stack: feature
   Layers:
   - layer1
     abc123 - commit1
     def456 - commit2

   - layer2
     ghi789 - commit3
     jkl012 - commit4

`

	if buf.String() != expectedOutput {
		t.Errorf("expected %q, got %q", expectedOutput, buf.String())
	}
}
