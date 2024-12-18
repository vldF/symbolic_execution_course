package testdata

func ComplexReal(a complex128) float64 {
	return real(a)
}

func ComplexImag(a complex128) float64 {
	return imag(a)
}

func ComplexId(a complex128) complex128 {
	return a
}

func BasicComplexOperations(a complex128, b complex128) complex128 {
	if real(a) > real(b) {
		return a + b
	} else if imag(a) > imag(b) {
		return a - b
	}
	return a * b
}

func ComplexMagnitude(a complex128) float64 {
	magnitude := real(a)*real(a) + imag(a)*imag(a)
	return magnitude
}

func ComplexComparison(a complex128, b complex128) int {
	magA := ComplexMagnitude(a)
	magB := ComplexMagnitude(b)

	if magA > magB {
		return 1
	} else if magA < magB {
		return -1
	}
	return 0
}

func ComplexOperations(a complex128, b complex128) complex128 {
	if real(a) == 0 && imag(a) == 0 {
		return b
	} else if real(b) == 0 && imag(b) == 0 {
		return a
	} else if real(a) > real(b) {
		return a / b
	}
	return a + b
}

func NestedComplexOperations(a complex128, b complex128) complex128 {
	if real(a) < 0 {
		if imag(a) < 0 {
			return a * b
		}
		return a + b
	}

	if imag(b) < 0 {
		return a - b
	}
	return a + b
}
