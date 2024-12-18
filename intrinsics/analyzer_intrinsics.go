package testdata

func Assume(predicate bool) {}

func MakeSymbolic[T any]() T {
	panic("shouldn't be called!")
}
