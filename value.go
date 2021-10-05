package jsonom

import (
	"encoding/json"
	"fmt"
	"io"
)

// A Value is one of: nil, string, bool, json.Number
type Value struct {
	value interface{}
}

func (v Value) isNode() {}

// NewValue returns a new value. The value must be convertible to one
// of: nil, string, bool, or json.Number
func NewValue(value interface{}) *Value {
	if value == nil {
		return &Value{}
	}
	switch value.(type) {
	case bool, string, json.Number:
		return &Value{value: value}
	case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint, float32, float64:
		b, _ := json.Marshal(value)
		return &Value{value: json.Number(b)}
	}
	panic(fmt.Sprint("Cannot decode value: %v", value))
}

// StringValue creates a new Value from a string
func StringValue(str string) *Value {
	return &Value{value: str}
}

// BoolValue creates a new value from a boolean value
func BoolValue(v bool) *Value {
	return &Value{value: v}
}

// NullValue returns a new null value
func NullValue() *Value {
	return &Value{}
}

func (v Value) String() string {
	if v.value == nil {
		return "null"
	}
	return fmt.Sprint(v.value)
}

// Value returns the value
func (v Value) Value() interface{} {
	return v.value
}

// Set sets the value
func (v *Value) Set(val interface{}) {
	v.value = val
}

// Marshal a value to interface{}
func (v Value) Marshal() interface{} {
	return v.value
}

// MarshalJSON allows using json.Marshal for a value node
func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Marshal())
}

// Encode a value
func (v Value) Encode(w io.Writer) error {
	data, err := json.Marshal(v.value)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
