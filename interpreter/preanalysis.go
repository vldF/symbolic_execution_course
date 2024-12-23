package interpreter

import "golang.org/x/tools/go/ssa"

func getAllBasicBlocks(f *ssa.Function) map[int]bool {
	result := make(map[int]bool)
	for _, block := range f.Blocks {
		result[block.Index] = true
	}

	return result
}
