package cmd

import (
	"bytes"
	"testing"

	"github.com/mattbearman/pancake/internal/git"
	"github.com/mattbearman/pancake/internal/stacks"
)

func TestShow(t *testing.T) {
	g := &git.Mock{
		Branch: "feature/layer2",
		LayerCommits: map[string][]string{
			"feature/layer1": {"abc123 - commit1", "def456 - commit2"},
			"feature/layer2": {"ghi789 - commit3", "jkl012 - commit4"},
		},
	}

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
	show(stackList, g, buf)

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
