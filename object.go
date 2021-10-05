package jsonom

import (
	"encoding/json"
	"io"
)

// KeyValue represents a JSON key-value pair in an Object
type KeyValue struct {
	key   string
	value Node
}

// NewKeyValue returns a new key-value pair
func NewKeyValue(key string, value Node) *KeyValue {
	if value == nil {
		value = &Value{}
	}
	return &KeyValue{key: key, value: value}
}

// Key returns the key of the key-value pair
func (k KeyValue) Key() string {
	return k.key
}

// Value returns the value of the key-value pair
func (k KeyValue) Value() Node {
	return k.value
}

// Set sets the value of the key-value pair
func (k *KeyValue) Set(value Node) {
	if value == nil {
		value = &Value{}
	}
	k.value = value
}

// Encode a key-value pair
func (k KeyValue) Encode(w io.Writer) error {
	data, err := json.Marshal(k.key)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	if _, err := w.Write([]byte{':'}); err != nil {
		return err
	}
	return k.value.Encode(w)
}

// Object represents a JSON object, an ordered set of key-value
// pairs. The zero-value of an Object is ready to use
type Object struct {
	kv    []*KeyValue
	kvMap map[string]int
}

func (o *Object) isNode() {}

// NewObject creates a new object from the given key-value pairs
func NewObject(kv ...*KeyValue) *Object {
	ret := &Object{}
	for _, x := range kv {
		ret.AddOrSet(x)
	}
	return ret
}

// Len returns the number of key-value pairs in the object
func (o *Object) Len() int {
	return len(o.kv)
}

// Get returns a key-value pair by its key. If the key does not exist, returns nil
func (o *Object) Get(key string) *KeyValue {
	if o.kvMap == nil {
		return nil
	}
	ix, ok := o.kvMap[key]
	if !ok {
		return nil
	}
	return o.kv[ix]
}

// Value returns the value of a key-value pair by its key. If the key
// does not exist, returns nil,false
func (o *Object) Value(key string) (Node, bool) {
	if o.kvMap == nil {
		return nil, false
	}
	k, ok := o.kvMap[key]
	if !ok {
		return nil, false
	}
	return o.kv[k].value, true
}

// N returns the n'th key-value pair by index
func (o *Object) N(index int) *KeyValue {
	return o.kv[index]
}

// AddOrSet adds a key-value pair if the key does not exist in this
// object, or replaces an existing key-value pair with the given
// one. Returns true if the key is added, false if it is updated
func (o *Object) AddOrSet(kv *KeyValue) bool {
	if o.kvMap == nil {
		o.kvMap = make(map[string]int)
	}
	old, exists := o.kvMap[kv.key]
	if exists {
		o.kv[old] = kv
		return false
	}
	ix := len(o.kv)
	o.kv = append(o.kv, kv)
	o.kvMap[kv.key] = ix
	return true
}

// Set the key to the given value. If the key already exists, it is
// replaced in place. Otherwise it is appended
func (o *Object) Set(key string, value Node) bool {
	return o.AddOrSet(NewKeyValue(key, value))
}

// Remove a key from the object. Returns true if the key was in the
// object, and removed
func (o *Object) Remove(key string) bool {
	if o.kvMap == nil {
		return false
	}
	ix, exists := o.kvMap[key]
	if !exists {
		return false
	}
	delete(o.kvMap, key)
	for k, v := range o.kvMap {
		if v > ix {
			o.kvMap[k] = ix - 1
		}
	}
	o.kv = append(o.kv[:ix], o.kv[ix+1:]...)
	return true
}

// Clear the contents of an object
func (o *Object) Clear() {
	o.kvMap = make(map[string]int)
	o.kv = make([]*KeyValue, 0)
}

// Marshal returns a map[string]interface{} for the object where each
// value is recursively marshaled. This operation loses the ordering
// of the object elements.
func (o *Object) Marshal() interface{} {
	ret := make(map[string]interface{}, len(o.kv))
	for _, x := range o.kv {
		var value interface{}
		if x.Value != nil {
			value = x.value.Marshal()
		}
		ret[x.key] = value
	}
	return ret
}

// MarshalJSON allows using json.Marshal for an object node
func (o Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Marshal())
}

// Encode a json object
func (o Object) Encode(w io.Writer) error {
	if _, err := w.Write([]byte{'{'}); err != nil {
		return err
	}
	for i, v := range o.kv {
		if i > 0 {
			if _, err := w.Write([]byte{','}); err != nil {
				return err
			}
		}
		if err := v.Encode(w); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{'}'}); err != nil {
		return err
	}
	return nil
}
