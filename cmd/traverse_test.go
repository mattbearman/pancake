package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/stretchr/testify/assert"
)

type gitForTraverseCmdMock struct {
	currentBranch string
}

func (g *gitForTraverseCmdMock) CurrentBranch() string {
	return g.currentBranch
}

func (g *gitForTraverseCmdMock) ChangeBranch(branchName string) {}

func TestUp(t *testing.T) {
	tests := []struct {
		name                  string
		gitForTraverseCmdMock GitForTraverseCmd
		expectedOutput        string
		expectedError         error
	}{
		{
			name:                  "moving up one layer",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "new-feature/part-2"},
			expectedOutput: `🥞 Working in git branch "new-feature/part-3"
`,
		},
		{
			name:                  "already at the top of the stack",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "new-feature/part-3"},
			expectedError:         fmt.Errorf("❌ you're already at the top of this stack"),
		},
		{
			name:                  "not in a stack",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "main"},
			expectedError:         fmt.Errorf(`❌ branch "main" is not part of a stack`),
		},
		{
			name:                  "unknown layer",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "new-feature/unknown"},
			expectedError:         fmt.Errorf(`❌ layer "unknown" is not part of the "new-feature" stack`),
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf(tt.name)

		t.Run(testname, func(t *testing.T) {
			buf := &bytes.Buffer{}

			stackFile := t.TempDir() + "/.pancake.json"
			stackList := *stacks.NewList(stackFile)
			stack := stackList.Add("new-feature", "main", "/")
			stack.AddLayer("part-1")
			stack.AddLayer("part-2")
			stack.AddLayer("part-3")

			err := up(stackList, tt.gitForTraverseCmdMock, buf)

			assert.Equal(t, tt.expectedOutput, buf.String())
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestDown(t *testing.T) {
	tests := []struct {
		name                  string
		gitForTraverseCmdMock GitForTraverseCmd
		expectedOutput        string
		expectedError         error
	}{
		{
			name:                  "moving down one layer",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "new-feature/part-3"},
			expectedOutput: `🥞 Working in git branch "new-feature/part-2"
`,
		},
		{
			name:                  "already at the bottom of the stack",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "new-feature/part-1"},
			expectedError:         fmt.Errorf("❌ you're already at the bottom of this stack"),
		},
		{
			name:                  "not in a stack",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "main"},
			expectedError:         fmt.Errorf(`❌ branch "main" is not part of a stack`),
		},
		{
			name:                  "unknown layer",
			gitForTraverseCmdMock: &gitForTraverseCmdMock{currentBranch: "new-feature/unknown"},
			expectedError:         fmt.Errorf(`❌ layer "unknown" is not part of the "new-feature" stack`),
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf(tt.name)

		t.Run(testname, func(t *testing.T) {
			buf := &bytes.Buffer{}

			stackFile := t.TempDir() + "/.pancake.json"
			stackList := *stacks.NewList(stackFile)
			stack := stackList.Add("new-feature", "main", "/")
			stack.AddLayer("part-1")
			stack.AddLayer("part-2")
			stack.AddLayer("part-3")

			err := down(stackList, tt.gitForTraverseCmdMock, buf)

			assert.Equal(t, tt.expectedOutput, buf.String())
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
