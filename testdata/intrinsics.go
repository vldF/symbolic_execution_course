package testdata

func SimpleAssume() int {
	val := MakeSymbolic[int]()
	Assume(val == 2)
	return val
}
