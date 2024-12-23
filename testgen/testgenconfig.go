package testgen

type PathSelectorMode int

const (
	Random PathSelectorMode = iota
	DFS    PathSelectorMode = iota
	NURS   PathSelectorMode = iota
)

type Config struct {
	MaxBasicBlocks   int
	PathSelectorMode PathSelectorMode
}
