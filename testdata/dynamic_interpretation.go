package testdata

func ImpossibleBranch(a int) int {
	if a > 0 {
		if a == 0 {
			return -1
		}
	} else if a < 0 {
		if a == 0 {
			return -1
		}
	}

	return a
}

func RepeatingConditions(a int) int {
	if a > 0 {
		return a
	} else if a > 0 {
		return -1
	}

	return a + 1
}

func ImpossibleCondition(a int) int {
	if a > 0 && a < 0 {
		return -1
	}

	return a
}
