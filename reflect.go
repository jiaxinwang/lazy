package lazy

import (
	"errors"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func clone(inter interface{}) interface{} {
	newInter := reflect.New(reflect.TypeOf(inter).Elem())

	val := reflect.ValueOf(inter).Elem()
	taggetVal := newInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		field := taggetVal.Field(i)
		field.Set(val.Field(i))
	}
	return newInter.Interface()
}

func deepCopy(src, dst interface{}) error {
	byt, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, dst)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func setJSONField(src interface{}, name string, v interface{}) error {
	byte1, err := json.Marshal(src)
	if err != nil {
		return err
	}
	var inter1 interface{}
	json.Unmarshal(byte1, &inter1)
	map1 := inter1.(map[string]interface{})
	byte2, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var inter2 interface{}
	json.Unmarshal(byte2, &inter2)
	map1[name] = inter2
	s, err := json.MarshalToString(map1)
	if err != nil {
		return err
	}
	json.UnmarshalFromString(s, src)
	return nil
}

func setField(src interface{}, name string, v interface{}) (err error) {
	ps := reflect.ValueOf(src)
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		f := s.FieldByName(name)
		if f.IsValid() {
			if f.CanSet() {
				if isNil(v) {
					f.Set(reflect.Zero(f.Type()))
				} else {
					f.Set(reflect.ValueOf(v))
				}
			} else {
				err = errors.New("can't set")
				return
			}
		} else {
			err = errors.New("not valid")
			return
		}
	} else {
		err = errors.New("wrong kind")
		return
	}
	return
}

func valueOfField(src interface{}, name string) (v interface{}, err error) {
	val := reflect.ValueOf(src).Elem()
	return val.FieldByName(name).Interface(), nil
}

func valueOfJSONKey(inter interface{}, key string) jsoniter.Any {
	if b, err := json.Marshal(inter); err == nil {
		return json.Get(b, key)
	}

	return nil
}

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
