package git

type Mock struct {
	Branch       string
	LayerCommits map[string][]string
}

func (g *Mock) CurrentBranch() string {
	return g.Branch
}

func (g *Mock) CommitsBetween(start, end string) []string {
	return g.LayerCommits[end]
}

func (g *Mock) RootDir() string {
	return "/"
}

func (g *Mock) NewBranch(name string) {
}
