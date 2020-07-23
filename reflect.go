package lazy

import (
	"reflect"
	"strings"
)

func valueOfTag(inter interface{}, tagName string) interface{} {
	model := reflect.ValueOf(inter)
	if model.Kind() == reflect.Ptr {
		model = model.Elem()
	}

	for i := 0; i < model.NumField(); i++ {
		fieldType := model.Type().Field(i)
		if name, _, _, _, err := disassembleTag(fieldType.Tag.Get("lazy")); err == nil && strings.EqualFold(name, tagName) {
			return model.Field(i).Interface()
		}
	}
	return nil
}

const (
	// ForeignOfModelName ...
	ForeignOfModelName = 0
	// ForeignOfModelID ...
	ForeignOfModelID = 1
	// ForeignOfModelForeignTable ...
	ForeignOfModelForeignTable = 2
	// ForeignOfModelForeignID ...
	ForeignOfModelForeignID = 3
)

func foreignOfModel(inter interface{}) [][4]string {
	ret := make([][4]string, 0)

	model := reflect.ValueOf(inter)
	if model.Kind() == reflect.Ptr {
		model = model.Elem()
	}
	for i := 0; i < model.NumField(); i++ {
		fieldType := model.Type().Field(i)
		if name, id, ft, fk, err := disassembleTag(fieldType.Tag.Get("lazy")); err == nil && len(ft) > 0 && len(fk) > 0 {
			ret = append(ret, [4]string{name, id, ft, fk})
		}
	}

	return ret
}
