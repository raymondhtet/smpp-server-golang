package helpers

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

func NewVisualizationObject[T any](headerData []byte, opts ...VisualizationObjectOptions[T]) *VisualizationObject[T] {
	vo := &VisualizationObject[T]{
		IsBigEndian: false,
		Bytes:       headerData,
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

func WithName[T any](name string) VisualizationObjectOptions[T] {
	return func(vo *VisualizationObject[T]) {
		vo.Name = name
	}
}

func (v *VisualizationObject[T]) Print() {
	print(v.Name, " -> Hex:", hex.EncodeToString(v.Bytes), " and Value:")

	value := v.GetValue()

	print(value, "\n")
}

func (v *VisualizationObject[T]) GetValue() T {

	switch any(v.Value).(type) {
	case uint32:
		v.Value = If[uint32](v.IsBigEndian, binary.BigEndian.Uint32(v.Bytes), binary.LittleEndian.Uint32(v.Bytes)).(T)
	default:
		v.Value = any(string(v.Bytes)).(T)
	}

	return v.Value
}

type Getter interface {
	Print()
	GetValue()
}

func AnyToBytes(data any) ([]byte, error) {
	// A type switch is used to check the actual concrete type of the interface.
	switch v := data.(type) {
	case []byte:
		// If it's already a byte slice, just return it.
		return v, nil
	case string:
		// If it's a string, convert it to a byte slice.
		return []byte(v), nil
	case int:
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(v))
		return b, nil
	case int8:
		return []byte{byte(v)}, nil
	case int16:
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(v))
		return b, nil
	case int32:
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(v))
		return b, nil
	case int64:
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(v))
		return b, nil
	case uint:
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(v))
		return b, nil
	case uint8:
		return []byte{v}, nil
	case uint16:
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, v)
		return b, nil
	case uint32:
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, v)
		return b, nil
	case uint64:
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, v)
		return b, nil

	default:
		// For any other type (like structs, maps, slices),
		// try to marshal it into JSON as a general-purpose fallback.
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("json.Marshal failed for type %T: %w", v, err)
		}
		return b, nil
	}
}
