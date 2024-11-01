package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattbearman/pancake/internal/stacks"
)

type gitForShowCmdMock struct {
	currentBranch string
}

func (g *gitForShowCmdMock) CurrentBranch() string {
	return g.currentBranch
}

func (g *gitForShowCmdMock) CommitsBetween(from string, to string) []string {
	switch to {
	case "feature/layer1":
		return []string{"abc123 - commit1", "def456 - commit2"}
	case "feature/layer2":
		return []string{"ghi789 - commit3", "jkl012 - commit4"}
	}

	return []string{}
}

func TestShow(t *testing.T) {
	tests := []struct {
		name              string
		gitForShowCmdMock GitForShowCmd
		expectedOutput    string
		expectedError     error
	}{
		{
			name: "show stack layers and commits",
			gitForShowCmdMock: &gitForShowCmdMock{
				currentBranch: "feature/layer2",
			},
			expectedOutput: `🥞 Current stack: feature
   Layers:
   - layer1
     abc123 - commit1
     def456 - commit2

   - layer2
     ghi789 - commit3
     jkl012 - commit4

`,
		},
		{
			name: "not in a stack",
			gitForShowCmdMock: &gitForShowCmdMock{
				currentBranch: "main",
			},
			expectedError: fmt.Errorf(`❌ branch "main" is not part of a stack`),
		},
		{
			name: "unknown layer",
			gitForShowCmdMock: &gitForShowCmdMock{
				currentBranch: "feature/unknown",
			},
			expectedError: fmt.Errorf(`❌ layer "unknown" is not part of the "feature" stack`),
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

	for _, tt := range tests {
		testname := fmt.Sprintf(tt.name)

		t.Run(testname, func(t *testing.T) {
			buf := &bytes.Buffer{}

			err := show(stackList, tt.gitForShowCmdMock, buf)

			assert.Equal(t, tt.expectedOutput, buf.String())
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
