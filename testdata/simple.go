package testdata

func IdInt(x int) int {
	return x
}

func IdFloat(x float64) float64 {
	return x
}

func SimpleExpressionInt(x int) int {
	t1 := x
	t2 := t1 + 1
	t3 := t2 + 2

	return t1 + t2 + t3
}
