package lazy

import (
	"reflect"
)

var registry = make(map[string]reflect.Type)

// Register ...
func Register(elems ...interface{}) {
	for _, elem := range elems {
		t := reflect.TypeOf(elem).Elem()
		registry[t.Name()] = t
	}
}

// NewStruct ...
func NewStruct(name string) (interface{}, bool) {
	elem, ok := registry[name]
	if !ok {
		return nil, false
	}
	return reflect.New(elem).Elem().Interface(), true
}

// NewStructSlice ...
func NewStructSlice(name string) (interface{}, bool) {
	elem, ok := registry[name]
	if !ok {
		return nil, false
	}
	return reflect.MakeSlice(reflect.SliceOf(elem), 0, 0).Interface(), true
}
