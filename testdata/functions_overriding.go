package testdata

import "math"

func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}

func math_Sqrt(x float64) float64 {
	result := MakeSymbolic[float64]()

	cond := result*result - x
	Assume(cond > -1e-3)
	Assume(cond < 1e-3)
	Assume(result >= 0)

	return result
}
