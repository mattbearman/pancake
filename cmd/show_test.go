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
		expectedExitCode  int
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
			expectedExitCode: 0,
		},
		{
			name: "not in a stack",
			gitForShowCmdMock: &gitForShowCmdMock{
				currentBranch: "main",
			},
			expectedOutput:   `❌ branch "main" is not part of a stack`,
			expectedExitCode: 1,
		},
		{
			name: "unknown layer",
			gitForShowCmdMock: &gitForShowCmdMock{
				currentBranch: "feature/unknown",
			},
			expectedOutput:   `❌ layer "unknown" is not part of the "feature" stack`,
			expectedExitCode: 1,
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
			expectedOutput := fmt.Sprintf("%s\n", tt.expectedOutput)

			exitCode := show(stackList, tt.gitForShowCmdMock, buf)

			assert.Equal(t, expectedOutput, buf.String())
			assert.Equal(t, tt.expectedExitCode, exitCode)
		})
	}
}
