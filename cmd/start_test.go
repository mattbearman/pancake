package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/mattbearman/pancake/internal/stacks"
)

type gitStarterMock struct{}

func (g *gitStarterMock) CurrentBranch() string {
	return "main"
}

func (g *gitStarterMock) NewBranch(branchName string) {}

func TestStart(t *testing.T) {
	tests := []struct {
		name           string
		gitStarterMock GitStarter
		args           []string
		separator      string
		expectedOutput string
		expectedJSON   string
	}{
		{
			name:           "start a new stack with default layer name",
			gitStarterMock: &gitStarterMock{},
			args:           []string{"new-feature"},
			separator:      "/",
			expectedOutput: `🥞 Created new stack "new-feature". Working in git branch "new-feature/part-1"`,
			expectedJSON:   `[{"name":"new-feature","baseBranch":"main","separator":"/","layers":["part-1"]}]`,
		},
		{
			name:           "start a new stack with specified layer name",
			gitStarterMock: &gitStarterMock{},
			args:           []string{"new-feature", "setup"},
			separator:      "/",
			expectedOutput: `🥞 Created new stack "new-feature". Working in git branch "new-feature/setup"`,
			expectedJSON:   `[{"name":"new-feature","baseBranch":"main","separator":"/","layers":["setup"]}]`,
		},
		{
			name:           "start a new stack with specified layer name and separator",
			gitStarterMock: &gitStarterMock{},
			args:           []string{"new-feature", "init"},
			separator:      "_",
			expectedOutput: `🥞 Created new stack "new-feature". Working in git branch "new-feature_init"`,
			expectedJSON:   `[{"name":"new-feature","baseBranch":"main","separator":"_","layers":["init"]}]`,
		},
	}

	buf := &bytes.Buffer{}

	for _, tt := range tests {
		testname := fmt.Sprintf(tt.name)

		t.Run(testname, func(t *testing.T) {
			stackFile := t.TempDir() + "/.pancake.json"
			stackList := *stacks.NewList(stackFile)

			separator = tt.separator

			start(stackList, tt.gitStarterMock, buf, tt.args...)

			expectedOutput := fmt.Sprintf("%s\n", tt.expectedOutput)

			if buf.String() != expectedOutput {
				t.Errorf("expected %q, got %q", expectedOutput, buf.String())
			}

			dat, err := os.ReadFile(stackFile)

			if err == nil {
				json := string(dat)

				if json != tt.expectedJSON {
					t.Errorf("expected %q, got %q", tt.expectedJSON, json)
				}
			} else {
				t.Error(err)
			}

			buf.Reset()
		})
	}
}
