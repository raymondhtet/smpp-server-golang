package helpers

import (
	"encoding/binary"
	"encoding/hex"
)

func If[T any](condition bool, conditionTrue T, conditionFalse T) any {
	if condition {
		return conditionTrue
	}
	return conditionFalse
}

type VisualizationObject[T any] struct {
	Bytes       []byte
	Value       T
	IsBigEndian bool
	Name        string
}

type VisualizationObjectOptions[T any] func(object *VisualizationObject[T])

func NewVisualizationObject[T any](headerData []byte, variableName string, opts ...VisualizationObjectOptions[T]) *VisualizationObject[T] {
	vo := &VisualizationObject[T]{
		IsBigEndian: false,
		Bytes:       headerData,
		Name:        variableName,
	}

	for _, opt := range opts {
		opt(vo)
	}

	return vo
}

func WithBigEndian[T any](isBigEndian bool) VisualizationObjectOptions[T] {
	return func(vo *VisualizationObject[T]) {
		vo.IsBigEndian = isBigEndian
	}
}
func (v *VisualizationObject[T]) Print() {
	print(v.Name, " -> Hex:", hex.EncodeToString(v.Bytes), " and Value:")

	switch any(v.Value).(type) {
	case uint32:
		v.Value = If[uint32](v.IsBigEndian, binary.BigEndian.Uint32(v.Bytes), binary.LittleEndian.Uint32(v.Bytes)).(T)
	default:
		v.Value = any(string(v.Bytes)).(T)
	}

	print(v.Value, "\n")
}

type Print interface {
	Print()
}
