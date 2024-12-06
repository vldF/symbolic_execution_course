package testdata

func ConstantLoop(a int) int {
	result := 0
	for i := 0; i < 10; i++ {
		result += a
	}

	if result == a {
		return 0
	}

	return result
}

func DynamicLoop(a int) int {
	result := 0
	for i := 0; i < a; i++ {
		result++
	}

	if result < 0 {
		return -1
	}

	return result
}
