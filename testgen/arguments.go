package testgen

import (
	"fmt"
	"strings"
)

type PrimitiveArgument struct {
	stringValue string
}

func (pa *PrimitiveArgument) String() string {
	return pa.stringValue
}

type PointerArgument struct {
	innerArgument fmt.Stringer
}

func (pa *PointerArgument) String() string {
	return "&" + pa.innerArgument.String()
}

type ArrayArgument struct {
	elementType string
	values      []fmt.Stringer
}

func (aa *ArrayArgument) String() string {
	argsAsString := make([]string, len(aa.values))
	for i, value := range aa.values {
		argsAsString[i] = value.String()
	}

	return "[]" + aa.elementType + "{" + strings.Join(argsAsString, ", ") + "}"
}

type StructArgument struct {
	name     string
	elements map[string]fmt.Stringer
}

func (s *StructArgument) String() string {
	var res strings.Builder
	res.WriteString(s.name)
	res.WriteString("{")

	for name, value := range s.elements {
		res.WriteString(name)
		res.WriteString(": ")
		res.WriteString(value.String())
		res.WriteString(",")
	}

	res.WriteString("}")

	return res.String()
}

type ComplexArgument struct {
	real string
	imag string
}

func (c *ComplexArgument) String() string {
	var res strings.Builder
	res.WriteString("complex")
	res.WriteString("(")
	res.WriteString(c.real)
	res.WriteString(", ")
	res.WriteString(c.imag)
	res.WriteString(")")

	return res.String()
}
