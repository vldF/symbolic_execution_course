package testdata

import "math/rand"

func InvokeExternal() int {
	a := rand.Int()
	b := rand.Int()
	c := rand.Int()

	if a+b+c == 1 {
		return 1
	}

	if a+b+c == 2 {
		return 2
	}

	if a+b+c == 3 {
		return 3
	}

	return -1
}
