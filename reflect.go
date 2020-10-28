package lazy

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var cacheStore *sync.Map

func init() {
	cacheStore = &sync.Map{}
}

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

var TimeReflectType = reflect.TypeOf(time.Time{})

// ParseField ...
func ParseField(fieldStruct reflect.StructField) *Field {
	field := &Field{
		Name:              fieldStruct.Name,
		FieldType:         fieldStruct.Type,
		IndirectFieldType: fieldStruct.Type,
		StructField:       fieldStruct,
		Creatable:         true,
		Updatable:         true,
		Readable:          true,
	}
	for field.IndirectFieldType.Kind() == reflect.Ptr {
		field.IndirectFieldType = field.IndirectFieldType.Elem()
	}

	// TODO: Creatable, Updatable, Readable

	// field.TagSettings

	fieldValue := reflect.New(field.IndirectFieldType)

	var getRealFieldValue func(reflect.Value)
	getRealFieldValue = func(v reflect.Value) {
		rv := reflect.Indirect(v)
		if rv.Kind() == reflect.Struct && !rv.Type().ConvertibleTo(TimeReflectType) {
			for i := 0; i < rv.Type().NumField(); i++ {
				newFieldType := rv.Type().Field(i).Type
				for newFieldType.Kind() == reflect.Ptr {
					newFieldType = newFieldType.Elem()
				}

				fieldValue = reflect.New(newFieldType)

				if rv.Type() != reflect.Indirect(fieldValue).Type() {
					getRealFieldValue(fieldValue)
				}

				if fieldValue.IsValid() {
					return
				}

				for key, value := range ParseTagSetting(field.IndirectFieldType.Field(i).Tag.Get("gorm"), ";") {
					if _, ok := field.TagSettings[key]; !ok {
						field.TagSettings[key] = value
					}
				}
			}
		}
	}

	getRealFieldValue(fieldValue)

	if _, ok := field.TagSettings["-"]; ok {
		field.Creatable = false
		field.Updatable = false
		field.Readable = false
		// field.DataType = ""
	}

	if v, ok := field.TagSettings["->"]; ok {
		field.Creatable = false
		field.Updatable = false
		if strings.ToLower(v) == "false" {
			field.Readable = false
		} else {
			field.Readable = true
		}
	}

	if v, ok := field.TagSettings["<-"]; ok {
		field.Creatable = true
		field.Updatable = true

		if v != "<-" {
			if !strings.Contains(v, "create") {
				field.Creatable = false
			}

			if !strings.Contains(v, "update") {
				field.Updatable = false
			}
		}
	}

	return field
}

// ParseDescribe ...
func ParseDescribe(inter interface{}) (*Describe, error) {
	if inter == nil {
		return nil, fmt.Errorf("%w: %+v", ErrNil, inter)
	}
	interType := reflect.ValueOf(inter).Type()
	for interType.Kind() == reflect.Slice || interType.Kind() == reflect.Array || interType.Kind() == reflect.Ptr {
		interType = interType.Elem()
	}

	if interType.Kind() != reflect.Struct {
		if interType.PkgPath() == "" {
			return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, inter)
		}
		return nil, fmt.Errorf("%w: %v.%v", ErrUnsupportedDataType, interType.PkgPath(), interType.Name())
	}

	if v, ok := cacheStore.Load(interType); ok {
		return v.(*Describe), nil
	}

	describe := &Describe{
		Name:         interType.Name(),
		ModelType:    interType,
		FieldsByName: map[string]*Field{},
	}

	// defer?

	for i := 0; i < interType.NumField(); i++ {
		if fieldStruct := interType.Field(i); ast.IsExported(fieldStruct.Name) {
			field := ParseField(fieldStruct)
			describe.Fields = append(describe.Fields, field)
		}
	}

	return describe, nil
}

// ParseTagSetting ...
func ParseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i
		if len(names[j]) > 0 {
			for {
				if names[j][len(names[j])-1] == '\\' {
					i++
					names[j] = names[j][0:len(names[j])-1] + sep + names[i]
					names[i] = ""
				} else {
					break
				}
			}
		}

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(strings.ToUpper(values[0]))

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = k
		}
	}

	return settings
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
