package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/mattbearman/pancake/internal/stacks"
	"github.com/stretchr/testify/assert"
)

type gitForNextCmdMock struct {
	currentBranch string
}

func (g *gitForNextCmdMock) CurrentBranch() string {
	return g.currentBranch
}

func (g *gitForNextCmdMock) NewBranch(branchName string) {}

func TestNext(t *testing.T) {
	tests := []struct {
		name              string
		gitForNextCmdMock GitForNextCmd
		args              []string
		separator         string
		expectedOutput    string
		expectedJSON      string
		expectedExitCode  int
	}{
		{
			name:              "start a new layer with default layer name",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature/part-2"},
			args:              []string{},
			separator:         "/",
			expectedOutput:    `🥞 Created new layer "part-3". Working in git branch "new-feature/part-3"`,
			expectedJSON:      `[{"name":"new-feature","baseBranch":"main","separator":"/","layers":["part-1","part-2","part-3"]}]`,
			expectedExitCode:  0,
		},
		{
			name:              "start a new layer with specified layer name",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature/part-2"},
			args:              []string{"setup"},
			separator:         "/",
			expectedOutput:    `🥞 Created new layer "setup". Working in git branch "new-feature/setup"`,
			expectedJSON:      `[{"name":"new-feature","baseBranch":"main","separator":"/","layers":["part-1","part-2","setup"]}]`,
			expectedExitCode:  0,
		},
		{
			name:              "start a new layer with specified layer name and separator",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature_part-2"},
			args:              []string{"init"},
			separator:         "_",
			expectedOutput:    `🥞 Created new layer "init". Working in git branch "new-feature_init"`,
			expectedJSON:      `[{"name":"new-feature","baseBranch":"main","separator":"_","layers":["part-1","part-2","init"]}]`,
			expectedExitCode:  0,
		},
		{
			name:              "not in a stack",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "main"},
			args:              []string{},
			separator:         "/",
			expectedOutput:    `❌ branch "main" is not part of a stack`,
			expectedExitCode:  1,
		},
		{
			name:              "unknown layer",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature/unknown"},
			args:              []string{},
			separator:         "/",
			expectedOutput:    `❌ layer "unknown" is not part of the "new-feature" stack`,
			expectedExitCode:  1,
		},
		// TODO: Add test for when not at the top of the stack
		// TODO: Add test for when layer name already exists in stack
	}

	for _, tt := range tests {
		testname := fmt.Sprintf(tt.name)

		t.Run(testname, func(t *testing.T) {
			buf := &bytes.Buffer{}

			stackFile := t.TempDir() + "/.pancake.json"
			stackList := *stacks.NewList(stackFile)
			stack := stackList.Add("new-feature", "main", tt.separator)
			stack.AddLayer("part-1")
			stack.AddLayer("part-2")

			separator = tt.separator

			expectedOutput := fmt.Sprintf("%s\n", tt.expectedOutput)

			exitCode := next(stackList, tt.gitForNextCmdMock, buf, tt.args...)

			assert.Equal(t, expectedOutput, buf.String())
			assert.Equal(t, tt.expectedExitCode, exitCode)

			dat, err := os.ReadFile(stackFile)

			if err == nil {
				json := string(dat)

				assert.Equal(t, tt.expectedJSON, json)
			} else if tt.expectedJSON != "" {
				t.Error(err)
			}
		})
	}
}
