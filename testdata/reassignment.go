package testdata

// always returns 1
func ArrayReassignment(arr []int) int {
	if arr[0] != 1 {
		arr[0] = 1
	}

	return arr[0]
}
