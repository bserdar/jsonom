package jsonom

import (
	"io"
)

// Node is a node in a JSON object model
type Node interface {
	// Marshal a node to one of: map[string]interface{}, []interface{},
	// or one of the JSON value types: nil, string, bool, or json.Number
	Marshal() interface{}
	// Encode the node into a JSON document
	Encode(io.Writer) error

	isNode()
}
