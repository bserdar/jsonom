package jsonom

import (
	"encoding/json"
	"io"
)

// Array represents a JSON array
type Array struct {
	nodes []Node
}

func (a *Array) isNode() {}

// NewArray returns an array with the given nodes
func NewArray(el ...Node) *Array {
	return &Array{nodes: el}
}

// Len returns the number of elements in the array
func (a *Array) Len() int {
	return len(a.nodes)
}

// N returns the n'th element of the array
func (a *Array) N(index int) Node {
	return a.nodes[index]
}

// Set the node at the given index
func (a *Array) Set(index int, value Node) {
	a.nodes[index] = value
}

// Append new values to the array
func (a *Array) Append(values ...Node) {
	a.nodes = append(a.nodes, values...)
}

// Remove the node at the given index
func (a *Array) Remove(index int) {
	a.nodes = append(a.nodes[:index], a.nodes[index+1:]...)
}

// Clear all nodes of an array
func (a *Array) Clear() {
	a.nodes = make([]Node, 0)
}

// Marshal returns a []interface{} for the array, where each element
// is recursively marshaled
func (a *Array) Marshal() interface{} {
	ret := make([]interface{}, 0, len(a.nodes))
	for _, x := range a.nodes {
		var value interface{}
		if x != nil {
			value = x.Marshal()
		}
		ret = append(ret, value)
	}
	return ret
}

// MarshalJSON allows using json.Marshal for an array node
func (a *Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Marshal())
}

// Encode array as JSON
func (a *Array) Encode(w io.Writer) error {
	if _, err := w.Write([]byte{'['}); err != nil {
		return err
	}
	for i, value := range a.nodes {
		if i > 0 {
			if _, err := w.Write([]byte{','}); err != nil {
				return err
			}
		}
		if err := value.Encode(w); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{']'}); err != nil {
		return err
	}
	return nil
}

// Call f for each index-value until it returns false
func (a *Array) Each(f func(int, Node) bool) bool {
	for index, item := range a.nodes {
		if !f(index, item) {
			return false
		}
	}
	return true
}
