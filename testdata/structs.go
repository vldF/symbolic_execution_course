package testdata

type Struct1 struct {
	IntField   int
	FloatField float64
}

func TestStruct(s Struct1) float64 {
	return float64(s.IntField) + s.FloatField
}
