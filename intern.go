package jsonom

// StringInterner is used to reduce the memory footprint of a json
// document by interning the keys, so the same copy of string is used
// throughout the file
type StringInterner interface {
	Intern(string) string
}

// MapInterner uses a map to keep single copies of strings. Empty
// value of DefaultInterner is ready to use
type MapInterner struct {
	strings map[string]string
}

// Intern updates the internal string table to include the given string
func (interner *MapInterner) Intern(s string) string {
	if interner.strings == nil {
		interner.strings = make(map[string]string, 128)
	}
	ret, ok := interner.strings[s]
	if ok {
		return ret
	}
	interner.strings[s] = s
	return s
}
