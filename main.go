package main

func main() {
	//	code := `
	//package main
	//
	//func integerOperations(a int, b int) int {
	//	if a > b {
	//		return a + b
	//	} else if a < b {
	//		return a - b
	//	} else {
	//		return a * b
	//	}
	//}
	//`

	code := `
package main

func integerOperations(a int, b int) int {
	if a > b {
		return a + b
	} else if a < b {
		return a - b
	} else {
		return a * b
	}
}
`
	BuildConstraints(code)
}
