package lazy

import (
	"reflect"
)

var registry = make(map[string]reflect.Type)

func register(elems ...interface{}) {
	for _, elem := range elems {
		t := reflect.TypeOf(elem).Elem()
		registry[t.Name()] = t
	}
}

func newStruct(name string) (interface{}, bool) {
	elem, ok := registry[name]
	if !ok {
		return nil, false
	}
	return reflect.New(elem).Elem().Interface(), true
}

func newStructSlice(name string) (interface{}, bool) {
	elem, ok := registry[name]
	if !ok {
		return nil, false
	}
	return reflect.MakeSlice(reflect.SliceOf(elem), 0, 0).Interface(), true
}
