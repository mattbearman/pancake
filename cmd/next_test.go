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
		expectedError     error
	}{
		{
			name:              "start a new layer with default layer name",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature/part-2"},
			separator:         "/",
			expectedOutput: `🥞 Created new layer "part-3". Working in git branch "new-feature/part-3"
`,
			expectedJSON: `[{"name":"new-feature","baseBranch":"main","separator":"/","layers":["part-1","part-2","part-3"]}]`,
		},
		{
			name:              "start a new layer with specified layer name",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature/part-2"},
			args:              []string{"setup"},
			separator:         "/",
			expectedOutput: `🥞 Created new layer "setup". Working in git branch "new-feature/setup"
`,
			expectedJSON: `[{"name":"new-feature","baseBranch":"main","separator":"/","layers":["part-1","part-2","setup"]}]`,
		},
		{
			name:              "start a new layer with specified layer name and separator",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature_part-2"},
			args:              []string{"init"},
			separator:         "_",
			expectedOutput: `🥞 Created new layer "init". Working in git branch "new-feature_init"
`,
			expectedJSON: `[{"name":"new-feature","baseBranch":"main","separator":"_","layers":["part-1","part-2","init"]}]`,
		},
		{
			name:              "not in a stack",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "main"},
			separator:         "/",
			expectedError:     fmt.Errorf(`❌ branch "main" is not part of a stack`),
		},
		{
			name:              "unknown layer",
			gitForNextCmdMock: &gitForNextCmdMock{currentBranch: "new-feature/unknown"},
			separator:         "/",
			expectedError:     fmt.Errorf(`❌ layer "unknown" is not part of the "new-feature" stack`),
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

			err := next(stackList, tt.gitForNextCmdMock, buf, tt.args...)

			assert.Equal(t, tt.expectedOutput, buf.String())
			assert.Equal(t, tt.expectedError, err)

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
