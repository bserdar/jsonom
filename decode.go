package jsonom

import (
	"bytes"
	"encoding/json"
	"io"
)

// Unmarshal the given byte slice to a json object model node
func Unmarshal(input []byte, interner StringInterner) (Node, error) {
	return UnmarshalReader(bytes.NewReader(input), interner)
}

// UnmarshalReader unmarshals the input to a json object model node
func UnmarshalReader(input io.Reader, interner StringInterner) (Node, error) {
	dec := json.NewDecoder(input)
	dec.UseNumber()
	return Decode(dec, interner)
}

// Decode a JSON object using the given decoder. Interner is optional,
// it will be used if given. If omitted, an internal temporary
// interner will be used.
func Decode(decoder *json.Decoder, interner StringInterner) (Node, error) {
	var ret Node

	if interner == nil {
		interner = &MapInterner{}
	}

	tok, err := decoder.Token()
	if err == io.EOF {
		return ret, nil
	}
	if err != nil {
		return nil, err
	}
	if delim, ok := tok.(json.Delim); ok {
		switch delim {
		case '{':
			ret, err = decodeObject(decoder, interner)
		case '[':
			ret, err = decodeArray(decoder, interner)
		default:
			err = &json.SyntaxError{Offset: decoder.InputOffset()}
		}
	} else {
		ret = decodeValue(tok)
	}
	return ret, err
}

func decodeObject(decoder *json.Decoder, interner StringInterner) (*Object, error) {
	ret := &Object{}
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			return ret, &json.SyntaxError{Offset: decoder.InputOffset()}
		}
		if err != nil {
			return ret, err
		}

		if delim, ok := tok.(json.Delim); ok {
			if delim == '}' {
				break
			}
			return ret, &json.SyntaxError{Offset: decoder.InputOffset()}
		}

		key, ok := tok.(string)
		if !ok {
			return ret, &json.SyntaxError{Offset: decoder.InputOffset()}
		}
		key = interner.Intern(key)

		value, err := Decode(decoder, interner)
		if err != nil {
			return ret, err
		}
		ret.AddOrSet(&KeyValue{key: key, value: value})
	}
	return ret, nil
}

func decodeElement(decoder *json.Decoder, interner StringInterner) (Node, bool, error) {
	var ret Node

	tok, err := decoder.Token()
	if err == io.EOF {
		return ret, false, &json.SyntaxError{Offset: decoder.InputOffset()}
	}
	if err != nil {
		return nil, false, err
	}
	if delim, ok := tok.(json.Delim); ok {
		switch delim {
		case '{':
			ret, err = decodeObject(decoder, interner)
		case '[':
			ret, err = decodeArray(decoder, interner)
		case ']':
			return ret, true, nil
		default:
			err = &json.SyntaxError{Offset: decoder.InputOffset()}
		}
	} else {
		ret = decodeValue(tok)
	}
	return ret, false, err
}

func decodeArray(decoder *json.Decoder, interner StringInterner) (*Array, error) {
	ret := &Array{}
	for {
		value, done, err := decodeElement(decoder, interner)
		if err != nil {
			return ret, err
		}
		if done {
			break
		}
		ret.Append(value)
	}
	return ret, nil
}

func decodeValue(tok json.Token) *Value {
	ret := &Value{}
	if tok == nil {
		return ret
	}
	switch val := tok.(type) {
	case bool:
		ret.value = val
	case json.Number:
		ret.value = val
	case string:
		ret.value = val
	case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint, float32, float64:
		b, _ := json.Marshal(val)
		ret.value = json.Number(b)
	}
	return ret
}
