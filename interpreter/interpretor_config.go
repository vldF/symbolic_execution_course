package interpreter

type PathSelectorMode int

const (
	Random PathSelectorMode = iota
	DFS    PathSelectorMode = iota
	NURS   PathSelectorMode = iota
)

type InterpreterConfig struct {
	PathSelectorMode PathSelectorMode
}
