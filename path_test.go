package jsonom

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPath(t *testing.T) {
	testCases := [][2]string{
		{`{"a":{"b":{"c":"str"}}}`, `["str"]`},
	}
	for _, testCase := range testCases {
		node, err := Unmarshal([]byte(testCase[0]), nil)
		if err != nil {
			t.Error(err)
			return
		}
		var v []any
		json.Unmarshal([]byte(testCase[1]), &v)
		result, err := ParseAndFind(node, "/a/b/c")
		if err != nil {
			t.Error(err)
			return
		}
		for i := range result {
			if v, ok := result[i].(*Value); ok {
				result[i] = v.value
			}
		}
		if !reflect.DeepEqual(result, v) {
			t.Errorf("Expected: %v got: %v", v, result)
		}
	}
}
