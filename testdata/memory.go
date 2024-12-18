package testdata

type Test struct {
	intField   int
	floatField float64
}

func StructAllocateStoreRead(a int) float64 {
	structure := &Test{
		intField:   0,
		floatField: 0,
	}

	structure.intField = a
	structure.floatField = float64(a)

	return float64(structure.intField) + structure.floatField
}

type Test2 struct {
	field *Test
}

func StructOfStructAllocateStoreRead(a int) float64 {
	structure1 := &Test{
		intField:   0,
		floatField: 0,
	}

	structure1.intField = a
	structure1.floatField = float64(a)

	structure2 := &Test2{
		field: structure1,
	}

	structure3 := &Test{
		intField:   0,
		floatField: 0,
	}

	structure2.field = structure3

	return float64(structure2.field.intField) + structure2.field.floatField
}

func ArrayAllocateStoreRead(a int) int {
	arr := make([]int, 3)
	arr[0] = a
	return arr[0]
}

func ArrayAllocateStoreReadDynamic(a int, idx int) int {
	arr := make([]int, idx+1)
	arr[idx] = a
	return arr[idx]
}

func ArrayAllocateStoreReadStore() int {
	arr := make([]int, 5)
	for i := 0; i < 5; i++ {
		arr[i] = 0
	}

	arr[2] = 1

	if arr[2] == 1 {
		arr[2] = 2
	}

	if arr[2] == 3 {
		return -1 // impossible
	}

	if arr[2] == 2 {
		arr[2] = 3
	}

	return arr[2]
}
