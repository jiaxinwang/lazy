package lazy

import (
	"fmt"
	"reflect"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
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

func setFieldWithJSONString(src interface{}, name string, v interface{}) error {
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

// Describe ...
type Describe struct {
	Name         string
	ModelType    reflect.Type
	Fields       []*Field
	FieldsByName map[string]*Field
}

// Field ...
type Field struct {
	Name              string
	FieldType         reflect.Type
	IndirectFieldType reflect.Type
	StructField       reflect.StructField
	Creatable         bool
	Updatable         bool
	Readable          bool
	Post              bool
	Put               bool
	Get               bool
	Patch             bool
	Tag               reflect.StructTag
	TagSettings       map[string]string
}

// TimeReflectType ...
var TimeReflectType = reflect.TypeOf(time.Time{})

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

func equalBuiltinInterface(a, b interface{}) bool {
	v1, err := builtinValue(a)
	if err != nil {
		return false
	}
	v2, err := builtinValue(b)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(v1, v2)

}

func assembleMany2Many(self []interface{}, join, foreign []map[string]interface{}, primaryKeyName, foreignKeyName, joinPrimaryKeyName, joinForeignKeyName, expJSONName string) ([]interface{}, error) {
	ret := make([]interface{}, len(self))
	for _, vSelf := range self {
		exp := make([]interface{}, 0)
		vPrimary, err := valueOfField(vSelf, primaryKeyName)
		if err != nil {
			return nil, err
		}
		for _, vJoin := range join {
			if vJoinPrimary, ok := vJoin[joinPrimaryKeyName]; ok {
				if equalBuiltinInterface(vPrimary, vJoinPrimary) {
					joinFKeyValue := vJoin[joinForeignKeyName]
					for _, vForeign := range foreign {
						if fKeyValue, ok := vForeign[foreignKeyName]; ok {
							if equalBuiltinInterface(fKeyValue, joinFKeyValue) {
								exp = append(exp, vForeign)
							}
						}
					}
				}
			}
		}
		if err := setFieldWithJSONString(vSelf, expJSONName, exp); err != nil {
			logrus.WithError(err).Error()
			return ret, err
		}
	}
	return ret, nil
}

func assembleHasMany(self []interface{}, foreign []map[string]interface{}, primaryKeyName, foreignKeyName, expJSONName string) ([]interface{}, error) {
	ret := make([]interface{}, len(self))
	for _, vSelf := range self {
		exp := make([]interface{}, 0)
		value, err := valueOfField(vSelf, primaryKeyName)
		if err != nil {
			return nil, err
		}

		for _, vForeign := range foreign {
			if sameKeyValue, ok := vForeign[foreignKeyName]; ok {
				v1, err := builtinValue(sameKeyValue)
				if err != nil {
					return nil, err
				}
				v2, err := builtinValue(value)
				if err != nil {
					return nil, err
				}
				if reflect.DeepEqual(v1, v2) {
					exp = append(exp, vForeign)
				}
			}
		}

		if err := setFieldWithJSONString(vSelf, expJSONName, exp); err != nil {
			logrus.WithError(err).Error()
			return ret, err
		}
	}
	return ret, nil
}

func assembleBelongTo(self []interface{}, foreign []map[string]interface{}, primaryKeyName, foreignKeyName, expJSONName string) ([]interface{}, error) {
	ret := make([]interface{}, len(self))
	for _, vSelf := range self {
		var exp interface{}
		value, err := valueOfField(vSelf, primaryKeyName)
		if err != nil {
			return nil, err
		}

		for _, vForeign := range foreign {
			if sameKeyValue, ok := vForeign[foreignKeyName]; ok {
				v1, err := builtinValue(sameKeyValue)
				if err != nil {
					return nil, err
				}
				v2, err := builtinValue(value)
				if err != nil {
					return nil, err
				}
				if reflect.DeepEqual(v1, v2) {
					exp = vForeign
				}
			}
		}

		if err := setFieldWithJSONString(vSelf, expJSONName, exp); err != nil {
			logrus.WithError(err).Error()
			return ret, err
		}
	}
	return ret, nil
}

func builtinValue(i interface{}) (o interface{}, err error) {

	value := reflect.ValueOf(i)

	switch value.Kind() {
	// case reflect.Int64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
	// 	o = int64(i.())
	// 	return
	// case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// 	o = i.(uint64)
	// 	return
	// case reflect.Uintptr:
	// 	// TODO:
	// 	return nil, fmt.Errorf("unsupported")

	case reflect.Uint:
		o = uint(i.(uint))
		return o, nil
	case reflect.Uint64:
		o = uint(i.(uint64))
		return o, nil
	case reflect.Uint32:
		o = uint(i.(uint32))
		return o, nil
	case reflect.Uint16:
		o = uint(i.(uint16))
		return o, nil
	case reflect.Uint8:
		o = uint(i.(uint8))
		return o, nil
	case reflect.Int:
		o = uint(i.(int))
		return o, nil
	case reflect.Int64:
		o = uint(i.(int64))
		return o, nil
	case reflect.Int32:
		o = uint(i.(int32))
		return o, nil
	case reflect.Int16:
		o = uint(i.(int16))
		return o, nil
	case reflect.Int8:
		o = uint(i.(int8))
		return o, nil
	case reflect.String:
		o = i.(string)
		return o, nil
	case reflect.Bool:
		// if kv, err := strconv.ParseBool(v); err == nil {
		// 	ret = kv
		// }
		return nil, fmt.Errorf("unsupported")
	default:
		// TODO:
		return nil, fmt.Errorf("unsupported")
	}
}
