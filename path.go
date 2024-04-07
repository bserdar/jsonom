package jsonom

import (
	"github.com/bserdar/gojsonpath"
)

// PathModel is the adapter for Node to use in JSON Path lookups
type PathModel struct {
	Node
}

func (m PathModel) Root() any { return m.Node }

func (PathModel) Type(node any) gojsonpath.NodeType {
	switch node.(type) {
	case *Object:
		return gojsonpath.ObjectNode
	case *Array:
		return gojsonpath.ArrayNode
	}
	return gojsonpath.ValueNode
}

func (PathModel) Len(in any) int {
	if arr, ok := in.(*Array); ok {
		return arr.Len()
	}
	return 0
}

func (PathModel) Key(node any, key string) (any, bool) {
	if obj, ok := node.(*Object); ok {
		return obj.Value(key)
	}
	return nil, false
}

func (PathModel) Keys(node any) []string {
	obj, ok := node.(*Object)
	if !ok {
		return nil
	}
	out := make([]string, 0, obj.Len())
	for i := 0; i < obj.Len(); i++ {
		out = append(out, obj.N(i).Key())
	}
	return out
}

func (PathModel) Elem(node any, index int) any {
	if arr, ok := node.(*Array); ok {
		return arr.N(index)
	}
	return nil
}

func (PathModel) Value(node any) any {
	v, ok := node.(*Value)
	if ok {
		return v.value
	}
	return nil
}

func Find(node Node, path gojsonpath.Path) ([]any, error) {
	return gojsonpath.Find(PathModel{node}, path)
}

// ParseAndFind parses the path and finds matching nodes in the doc
func ParseAndFind(node Node, path string) ([]any, error) {
	return gojsonpath.ParseAndFind(PathModel{node}, path)
}

// Search iterates all document nodes depth-first, and calls `capture`
// for those document nodes that `path` matches.
func Search(node Node, path gojsonpath.Path, capture func(gojsonpath.DocPath)) error {
	return gojsonpath.Search(PathModel{node}, path, capture)
}
