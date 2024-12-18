package interpreter

import "go/types"

type Mocker struct {
	history map[string][]MockDescriptor
}

type MockDescriptor struct {
	Memory *Memory
	Value  Value
	Type   types.Type
}

func (m *Mocker) Copy() *Mocker {
	newHistory := make(map[string][]MockDescriptor)

	for k, v := range m.history {
		newHistory[k] = v
	}

	return &Mocker{
		history: newHistory,
	}
}

func (m *Mocker) Add(
	funcName string,
	value Value,
	mem *Memory,
	typ types.Type,
) {
	newDescriptor := MockDescriptor{
		Memory: mem,
		Value:  value,
		Type:   typ,
	}
	if stored, ok := m.history[funcName]; ok {
		m.history[funcName] = append(stored, newDescriptor)
	} else {
		m.history[funcName] = []MockDescriptor{newDescriptor}
	}
}

func (m *Mocker) GetAll() map[string][]MockDescriptor {
	return m.history
}
