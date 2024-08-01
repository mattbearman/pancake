package stacks

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/mattbearman/pancake/internal/git"
)

var pancakeFile string = filepath.Join(git.WorkingDir(), ".pancake.json")

func Load() {
	dat, err := os.ReadFile(pancakeFile)

	if err != nil {
		// Swallow file not found error, as we will initialize an empty slice of stacks
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	} else {
		json.Unmarshal([]byte(dat), &allStacks)
	}
}

func Save() {
	dat, err := json.Marshal(allStacks)

	if err != nil {
		panic(err)
	}

	os.WriteFile(pancakeFile, dat, 0644)
}
