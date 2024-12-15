package testdata

func inc(a int) int {
	return a + 1
}

func Twice(a int) int {
	return inc(inc(a))
}

func incComplex(a complex128) complex128 {
	return a + complex(1, 1)
}

func TwiceComplex(a complex128) complex128 {
	return incComplex(incComplex(a))
}

type IncStruct struct {
	field int
}

func incStruct(val *IncStruct) *IncStruct {
	val.field++
	return val
}

func TwiceStruct(a int) int {
	t := &IncStruct{field: a}
	return incStruct(t).field
}

func AddRecursive(a int, b int) int {
	if b == 0 {
		return a
	}

	return AddRecursive(a+1, b-1)
}

func Fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return Fib(n-1) + Fib(n-2)
	}
}
