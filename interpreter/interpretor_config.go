package interpreter

type PathSelectorMode int

const (
	Random PathSelectorMode = iota
	DFS    PathSelectorMode = iota
	NURS   PathSelectorMode = iota
)

type InterpretationMode int

const (
	Execution            InterpretationMode = iota
	CoverageMaximization InterpretationMode = iota
)

type InterpreterConfig struct {
	PathSelectorMode PathSelectorMode
	MainPackage      string
	Mode             InterpretationMode
}
