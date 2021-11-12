[![GoDoc](https://godoc.org/github.com/bserdar/jsonom?status.svg)](https://godoc.org/github.com/bserdar/jsonom)
[![Go Report Card](https://goreportcard.com/badge/github.com/bserdar/jsonom)](https://goreportcard.com/report/github.com/bserdar/jsonom)

# JSON Object Model

This is a Go package to process JSON documents with an API similar to
XML DOM API. It represents the JSON document as a tree of Nodes.  The
basic JSON types are:

  * Object is a key-value pair, where value is another node,
  * Array is an ordered value list,
  * Value is a JSON value node containing null, string, boolean, or a
    number
    
When a JSON document is unmarshaled, primitive values are converted to
the following Go types:
  
  * JSON number values are converted json.Number
  * JSON strings are converted to Go strings
  * JSON boolean values are converted to Go boolean
  * Null is converted to nil

This library preserves the ordering of keys in a JSON object.

## Key Interning

In JSON documents, the objects keys tend to repeat. For large JSON
documents, it makes sense to "intern" these keys in a hashmap so a
single string copy is used for all instances of a string value. 

JSON unmarshaling functions get an `Interner` argument. This can be
used to keep a common interner when processing multiple documents. If
a `nil` interner is passed to these functions, an new interner will be
used to store the keys of that document only.

```
func main() {
      input:=`{"key":"value", "arr": [ 1,2 ]}`
      j, err:=jsonom.Unmarshal([]byte(input),nil)
      if err!=nil {
         panic(err)
      }
      obj:=j.(*jsonom.Object)
      v,_:=obj.Value("key")
      fmt.Println(v.(*jsonom.Value).Value())
      a,_:=obj.Value("arr")
      arr:=a.(*jsonom.Array)
      fmt.Println(arr.Len())
      arr.Append(jsonom.NewValue(3))
      fmt.Println(arr.Len())
      j.Encode(os.Stdout)
}
```
