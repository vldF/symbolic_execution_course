package testdata

func ShortOverflow(x int, y int) int {
	if y > 10 || y <= 0 {
		return 0
	}

	return int(int8(x) + int8(y))
}

func OverflowInLoop(x int) int {
	var res int8 = 120

	for i := 1; i <= x; i++ {
		res++
	}

	return int(res)
}
